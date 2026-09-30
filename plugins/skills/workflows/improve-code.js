export const meta = {
  name: 'improve-code',
  description: 'Review code through independent one-lens reviewers, triage, fix act-on findings one commit each, re-review until nothing is left to act on, then tidy the history',
  whenToUse: 'Review and improve existing code starting from a set of paths or a revision range. A range is reviewed as one set: every file it changes, in its current state. Args: {paths?: string[], revs?: string, rounds?: number (default 3), check?: string, model?: string, models?: {scan?, review?, triage?, fix?, tidy?}, efforts?: {same keys}}, or a plain string naming the scope. models and efforts set one stage; model sets every stage without its own entry. Scan (VCS, check command, conventions, docs, layout; facts only) defaults to haiku at low effort and its brief is given to every later stage; every other stage inherits the session model.',
  phases: [
    { title: 'Scan', detail: 'facts about the project: VCS, check command, conventions, docs, layout' },
    { title: 'Review', detail: 'one fresh reviewer per lens over the whole scope' },
    { title: 'Triage', detail: 'act-on, consider, or dismissed' },
    { title: 'Fix', detail: 'one fixer per round, one commit per finding' },
    { title: 'Tidy', detail: 'absorb or fold fix commits, re-run the check on each' },
  ],
}

// Lens text lives in the plugin's agent files; `reach` is how far past the
// scope its reviewer looks.
const lenses = {
  structural: {
    agentType: 'sebnow:reviewer-structural',
    reach: 'codebase',
  },
  general: {
    agentType: 'sebnow:reviewer-general',
    reach: 'neighbourhood',
  },
  tests: {
    agentType: 'sebnow:reviewer-tests',
    reach: 'unit',
  },
  duplication: {
    agentType: 'sebnow:reviewer-duplication',
    reach: 'codebase',
  },
  writing: {
    agentType: 'sebnow:reviewer-writing',
    reach: 'unit',
  },
}
const lensKeys = Object.keys(lenses)

const reachText = {
  unit: 'the scope\'s files only; when the scope holds tests, also the code they exercise',
  neighbourhood: 'the scope\'s files, their direct callers and callees, and the package or component that owns them',
  codebase: 'the whole repository; search it',
}

const evidenceKinds = ['traced the path', 'ran it', 'reasoned', 'hypothetical']
const vcsCommand = 'test -d .jj && echo jj || (git rev-parse --is-inside-work-tree >/dev/null 2>&1 && echo git || echo none)'

const scanSchema = {
  type: 'object',
  required: ['vcs', 'check', 'layout', 'conventions', 'docs', 'has_tests'],
  properties: {
    vcs: { enum: ['jj', 'git', 'none'] },
    check: { type: 'string' },
    layout: { type: 'string' },
    conventions: { type: 'string' },
    docs: { type: 'array', items: { type: 'string' } },
    has_tests: { type: 'boolean' },
  },
}

const findingsSchema = {
  type: 'object',
  required: ['findings', 'read_outside_scope'],
  properties: {
    read_outside_scope: { type: 'array', items: { type: 'string' } },
    findings: {
      type: 'array',
      items: {
        type: 'object',
        required: ['locations', 'problem', 'consequence', 'evidence', 'evidence_detail', 'severity', 'fix'],
        properties: {
          locations: { type: 'array', minItems: 1, items: { type: 'string', pattern: '^[^:\\s]+:\\d+$' } },
          problem: { type: 'string' },
          consequence: { type: 'string' },
          evidence: { enum: evidenceKinds },
          evidence_detail: { type: 'string' },
          severity: { enum: ['high', 'medium', 'low'] },
          fix: { type: 'string' },
        },
      },
    },
  },
}

const triageSchema = {
  type: 'object',
  required: ['decisions'],
  properties: {
    decisions: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'bucket', 'reason', 'recurs'],
        properties: {
          id: { type: 'string' },
          bucket: { enum: ['act-on', 'consider', 'dismissed'] },
          reason: { type: 'string' },
          recurs: { type: 'string' },
        },
      },
    },
  },
}

const fixSchema = {
  type: 'object',
  required: ['results'],
  properties: {
    results: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'status', 'commit', 'files_changed', 'check_command', 'check_output', 'note'],
        properties: {
          id: { type: 'string' },
          status: { enum: ['fixed', 'reverted'] },
          commit: { type: 'string' },
          files_changed: { type: 'array', items: { type: 'string' } },
          check_command: { type: 'string' },
          check_output: { type: 'string' },
          note: { type: 'string' },
        },
      },
    },
  },
}

const tidySchema = {
  type: 'object',
  required: ['commits', 'unabsorbed'],
  properties: {
    commits: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'message', 'check', 'check_output'],
        properties: {
          id: { type: 'string' },
          message: { type: 'string' },
          check: { enum: ['pass', 'fail'] },
          check_output: { type: 'string' },
        },
      },
    },
    unabsorbed: {
      type: 'array',
      items: {
        type: 'object',
        required: ['what', 'why'],
        properties: { what: { type: 'string' }, why: { type: 'string' } },
      },
    },
  },
}

function parseArgs(raw) {
  if (raw && typeof raw === 'object') return raw
  if (typeof raw !== 'string' || !raw.trim()) return {}
  try {
    const parsed = JSON.parse(raw)
    if (parsed && typeof parsed === 'object') return parsed
  } catch (e) {}
  return { scope: raw.trim() }
}

const input = parseArgs(args)
const scopeLines = []
if (Array.isArray(input.paths) && input.paths.length) scopeLines.push('Paths: ' + input.paths.join(', '))
if (input.revs) scopeLines.push('Revision range: ' + input.revs + '. Review the whole range as one set: every file it changes, in its current state. Do not review commit by commit.')
if (input.scope) scopeLines.push('Scope: ' + input.scope)
if (!scopeLines.length) {
  return { error: 'No scope given. Pass args {paths: [...]} or {revs: "<range>"}, or a string naming the scope.' }
}
const scope = scopeLines.join('\n')
const maxRounds = Number.isInteger(input.rounds) && input.rounds > 0 ? input.rounds : 3

// Scan gathers facts and judges nothing, so it defaults to a cheap tier.
// Every other stage judges or changes code and inherits the session model
// unless args say otherwise.
const defaultModels = { scan: 'haiku' }
const defaultEfforts = { scan: 'low' }
function stageOpts(stage, opts) {
  const model = (input.models || {})[stage] || input.model || defaultModels[stage]
  const effort = (input.efforts || {})[stage] || defaultEfforts[stage]
  return { ...opts, ...(model ? { model } : {}), ...(effort ? { effort } : {}) }
}

phase('Scan')
const scan = await agent(
  'Gather facts about this project for reviewers who have not seen it. Report what you observe; do not judge it.\n\n' +
  scope + '\n\n' +
  '- vcs: run exactly this command at the repository root and report its output verbatim:\n  ' + vcsCommand + '\n' +
  '- check: the one shell command that builds the project and runs its tests, found from the project files.\n' +
  '- layout: the components or packages and what each owns, in a few lines, with the scope placed among them.\n' +
  '- conventions: how the code is written, as observed in the files around the scope and in any style or contributor docs: naming, file and package layout, error handling, logging, test style, comment style, commit message style from the recent history.\n' +
  '- docs: paths of the files that state rules or context for this code (agent instructions, README, CONTRIBUTING, style guides, architecture notes), each with a one-line note of what it covers.\n' +
  '- has_tests: whether the scope contains automated tests.\n' +
  'Quote a rule rather than paraphrasing it. Do not edit anything.',
  stageOpts('scan', { label: 'scan', phase: 'Scan', schema: scanSchema }),
)
if (!scan) return { error: 'The scan agent returned nothing.' }
const check = input.check || scan.check
const vcs = scan.vcs
log('vcs: ' + vcs + '; check: ' + check)
const projectBrief =
  'PROJECT (gathered by a scan; facts, not judgments):\n' +
  'Layout: ' + scan.layout + '\n' +
  'Conventions: ' + scan.conventions + '\n' +
  'Docs worth reading: ' + (scan.docs.length ? scan.docs.join('; ') : '(none found)') + '\n\n'

const reviewPrompt = (lens, focus) =>
  'Review this scope through your lens.\n\n' + scope + '\n' +
  (focus.length ? 'Files changed by fixes since the last review; start from these and what they touch: ' + focus.join(', ') + '\n' : '') +
  '\n' + projectBrief +
  'Skip generated, vendored, and non-code files. Cover the whole scope; decide yourself what to read together. ' +
  'The scope is where you start, not a boundary. Your reach: ' + lenses[lens].reach + ', meaning ' + reachText[lenses[lens].reach] + '. ' +
  'A finding may point anywhere in the repository. List every file outside the scope that you read in read_outside_scope.\n\n' +
  'Report only problems your lens covers. For each finding give:\n' +
  '- locations: one or more path:line entries, one line number each\n' +
  '- problem: what is wrong\n' +
  '- consequence: what concretely goes wrong because of it, and when\n' +
  '- evidence: "traced the path" if you followed the code to the consequence, "ran it" if a command you ran shows it, ' +
  '"reasoned" if you argued it from reading without tracing, "hypothetical" if it depends on conditions you did not confirm\n' +
  '- evidence_detail: the path you traced, the command and its output, or the argument\n' +
  '- severity: high, medium, or low\n' +
  '- fix: a sketch of the change, naming every file it touches\n' +
  'Return an empty list when you find nothing; do not pad it.'

const triagePrompt = (findings, fixedSummary) =>
  'Triage these code review findings. Read the code to check each one; do not take the reviewer\'s word for it.\n\n' + projectBrief +
  'Buckets:\n' +
  '- act-on: the problem is real, the consequence follows from the evidence, and the fix needs no design decision from a human.\n' +
  '- consider: plausible or worth a human look, but not established, or the fix needs a design decision.\n' +
  '- dismissed: wrong, unreachable, a style preference, or a duplicate of another finding here (name its id).\n' +
  'A hypothetical finding is act-on only if you confirm it yourself. Give every finding exactly one decision with a reason.\n\n' +
  'Set recurs to the id of a PREVIOUSLY FIXED finding when this finding has the same locations and the same substance; otherwise set it to an empty string.\n\n' +
  'FINDINGS:\n' + JSON.stringify(findings, null, 1) + '\n\n' +
  'PREVIOUSLY FIXED:\n' + (fixedSummary.length ? JSON.stringify(fixedSummary, null, 1) : '(none)') + '\n\n' +
  'Do not edit anything.'

const fixPrompt = findings =>
  'Fix these review findings, in order. ' +
  'Fix each at its cause, wherever in the repository that is, even outside the reviewed scope. Follow the project conventions below.\n\n' + projectBrief +
  'VCS: ' + vcs + '. Check command: ' + check + '\n\n' +
  'The working copy may hold changes that are not yours: never commit, squash, or revert them. Commit only files you changed, by naming their paths.\n\n' +
  'Run the check once before editing. If it already fails, change nothing and report every finding as reverted with that output.\n\n' +
  'For each finding, in order:\n' +
  '1. Make the smallest change that removes the problem at its cause. Fix nothing that is not listed.\n' +
  '2. Run the check command. Report the exact command and the end of its output.\n' +
  '3. If it passes, commit that change alone, following the commit skill. The body names the finding: its id, locations, and problem. ' +
  'Report the commit id (jj change id or git short hash) and the files you changed.\n' +
  '4. If you cannot make the check pass, revert this finding\'s change, confirm the check passes again, and report status reverted with the failing output.\n\n' +
  'FINDINGS:\n' + JSON.stringify(findings, null, 1)

const roundRecords = []
const fixed = []
const escalated = []
let focus = []
const lensList = lensKeys.filter(k => k !== 'tests' || scan.has_tests)
if (lensList.length < lensKeys.length) log('No tests in scope; the tests lens is skipped')
let stopReason = 'round limit (' + maxRounds + ') reached'

for (let r = 1; r <= maxRounds; r++) {
  const record = { round: r, focus, reviews: [], review_failures: [], act_on: [], consider: [], dismissed: [], fixes: [] }
  roundRecords.push(record)

  const reviews = await pipeline(lensList, lens =>
    agent(reviewPrompt(lens, focus), stageOpts('review', {
      label: 'r' + r + ' ' + lens,
      phase: 'Review',
      schema: findingsSchema,
      agentType: lenses[lens].agentType,
    })),
  )

  const findings = []
  reviews.forEach((res, i) => {
    const lens = lensList[i]
    if (!res) {
      record.review_failures.push(lens)
      return
    }
    record.reviews.push({ lens, findings: res.findings.length, read_outside_scope: res.read_outside_scope })
    for (const f of res.findings) {
      findings.push({ id: 'r' + r + '-' + (findings.length + 1), lens, ...f })
    }
  })
  if (record.review_failures.length) log('Round ' + r + ' reviews that returned nothing: ' + record.review_failures.join(', '))
  log('Round ' + r + ': ' + findings.length + ' findings')

  if (!findings.length) {
    stopReason = 'no findings in round ' + r
    break
  }

  const fixedSummary = fixed.map(f => ({ id: f.id, locations: f.locations, problem: f.problem }))
  const triage = await agent(triagePrompt(findings, fixedSummary), stageOpts('triage', { label: 'r' + r + ' triage', phase: 'Triage', schema: triageSchema }))
  if (!triage) {
    record.consider = findings.map(f => ({ ...f, reason: 'triage returned nothing' }))
    stopReason = 'triage failed in round ' + r
    break
  }

  const decisions = new Map(triage.decisions.map(d => [d.id, d]))
  for (const f of findings) {
    const d = decisions.get(f.id) || { bucket: 'consider', reason: 'triage gave no decision', recurs: '' }
    const entry = { ...f, reason: d.reason }
    const prior = d.recurs ? fixed.find(x => x.id === d.recurs) : null
    if (d.bucket === 'act-on' && prior) escalated.push({ finding: entry, repeats: prior })
    else if (d.bucket === 'act-on') record.act_on.push(entry)
    else if (d.bucket === 'dismissed') record.dismissed.push(entry)
    else record.consider.push(entry)
  }
  log('Round ' + r + ' triage: ' + record.act_on.length + ' act-on, ' + record.consider.length + ' consider, ' + record.dismissed.length + ' dismissed, ' + escalated.length + ' recurring')

  if (escalated.length) {
    stopReason = 'a fixed finding recurred in round ' + r + '; escalated, not retried'
    break
  }
  if (!record.act_on.length) {
    stopReason = 'no act-on findings in round ' + r
    break
  }

  // One fixer per round: fixes share one working copy and one check.
  const changedFiles = new Set()
  const res = await agent(fixPrompt(record.act_on), stageOpts('fix', { label: 'r' + r + ' fix', phase: 'Fix', schema: fixSchema }))
  const results = res ? res.results : []
  for (const f of record.act_on) {
    const out = results.find(x => x.id === f.id) || {
      status: 'unreported', commit: '', files_changed: [], check_command: '', check_output: '',
      note: res ? 'fixer did not report this finding' : 'fixer returned nothing',
    }
    record.fixes.push({ id: f.id, locations: f.locations, ...out })
    if (out.status === 'fixed') {
      fixed.push({ ...f, commit: out.commit, files_changed: out.files_changed, round: r })
      out.files_changed.forEach(x => changedFiles.add(x))
    }
  }

  focus = [...changedFiles]
  if (!focus.length) {
    stopReason = 'no fix landed in round ' + r
    break
  }
}

let tidy = null
if (fixed.length) {
  phase('Tidy')
  const fixCommits = fixed.map(f => ({ id: f.id, round: f.round, commit: f.commit, locations: f.locations, problem: f.problem, files_changed: f.files_changed }))
  const range = input.revs
    ? 'The reviewed range is ' + input.revs + '; the fix commits sit on top of it.'
    : 'No revision range was reviewed; the range is the fix commits themselves.'
  tidy = await agent(
    'Tidy the history of these fix commits. Load and follow the commit skill, and the jujutsu skill if the VCS is jj. Do not change code.\n\n' + projectBrief +
    'VCS: ' + vcs + '. Check command: ' + check + '\n' +
    range + '\n\n' +
    'Rules:\n' +
    '1. For each fix, ask: does it belong to a commit that already exists in the range? ' +
    'If yes, absorb it into that commit. If it is separable, keep it as its own commit. ' +
    'A later-round fix of the same finding (same locations and substance) belongs to the earlier fix commit.\n' +
    '2. Absorbing. jj: jj absorb --from <fix> --into <range>; it moves each hunk into the commit that last touched those lines and leaves ambiguous hunks behind. ' +
    'A fully absorbed commit that has a description is not abandoned for you: abandon it once it is empty. ' +
    'git: git commit --fixup=<target>, then GIT_SEQUENCE_EDITOR=true git rebase -i --autosquash (add --autostash if the working copy has changes); use git absorb instead if it is installed.\n' +
    '3. Order the commits that remain so a reviewer can replay them as an argument: ' +
    'a deletion before the reshape it makes room for, then any follow-on cleanup; a scaffold before the feature that uses it; a test before the fix it pins. ' +
    'Every commit must land on its own: if a test fails without its fix, keep the two in one commit rather than leave a red commit.\n' +
    '4. Keep commits small: one logical change each, never a batch of fixes.\n' +
    '5. Messages follow the project\'s commit style, as the commit skill detects it. ' +
    'Describe the change and why, not the finding id. A body does not restate its subject.\n\n' +
    'The working copy may hold changes that are not yours: never commit, squash, or revert them. ' +
    'Note the working-copy change before you start and return to it at the end.\n\n' +
    'After rewriting, run the check on every rewritten commit. ' +
    'git: git rebase -x "' + check + '" over the rewritten range. ' +
    'jj: for each revision in the rewritten range, jj new <rev>, run the check, then jj abandon the empty working-copy commit.\n\n' +
    'Report every final commit in order with its id, message first line, and check result. ' +
    'A commit that fails its check is reported as fail; do not fix it. ' +
    'Report every fix or hunk you could not absorb or fold, and why.\n\n' +
    'FIX COMMITS:\n' + JSON.stringify(fixCommits, null, 1),
    stageOpts('tidy', { label: 'tidy', phase: 'Tidy', schema: tidySchema }),
  )
  if (!tidy) log('The tidy agent returned nothing; fix commits are as the fixers left them.')
}

return {
  stop_reason: stopReason,
  scope,
  check,
  vcs,
  lenses: lensList,
  rounds: roundRecords,
  fixed: fixed.map(f => ({ id: f.id, locations: f.locations, problem: f.problem, commit: f.commit, round: f.round })),
  escalated,
  tidy,
}

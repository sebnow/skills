package ccvr

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Clause keys name the clauses of CCVR 2021 that leave a choice to the
// operator. The key is what `tidewright elect` takes (REQ-0139).
const (
	ClauseSurveyDeferral   = "survey-deferral"
	ClauseLaidUpSeason     = "laid-up-season"
	ClauseLevyInstalments  = "levy-instalments"
	ClauseOverhaulRecord   = "overhaul-record"
	ClauseCorrectionMethod = "correction-method"
	ClauseSeasonStart      = "season-start"
	ClauseHullClass        = "hull-class"
	ClauseTenderAsVessel   = "tender-as-vessel"
)

// Input keys name the inputs a rule depends on that nothing records yet.
const (
	InputHomePort           = "home-port"
	InputForHire            = "for-hire"
	InputOffListYardHours   = "off-list-yard-hours"
	InputCarriedEngineHours = "carried-engine-hours"
)

// Clause is a clause of the Regulations that offers the operator a choice.
type Clause struct {
	Key     string
	Label   string
	Options []Option
}

// Option is one choice a clause offers.
type Option struct {
	Key string
}

// ClauseByKey returns the clause with the given key.
func ClauseByKey(key string) (Clause, bool) {
	for _, c := range registry() {
		if c.Key == key {
			return c, true
		}
	}
	return Clause{}, false
}

// Offers reports whether the clause offers the option.
func (c Clause) Offers(option string) bool {
	for _, o := range c.Options {
		if o.Key == option {
			return true
		}
	}
	return false
}

// OptionKeys lists the clause's options in order.
func (c Clause) OptionKeys() []string {
	keys := make([]string, len(c.Options))
	for i, o := range c.Options {
		keys[i] = o.Key
	}
	return keys
}

// EngineClass is an engine class in Schedule 2.
type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
	EngineClassD EngineClass = "D"
)

// overhaulThreshold is Schedule 2: hours between overhauls per engine class.
var overhaulThreshold = map[EngineClass]float64{
	EngineClassA: 1500,
	EngineClassB: 2000,
	EngineClassC: 6000,
	EngineClassD: 8000,
}

// engineClasses maps the engine models the importers report to their
// Schedule 2 class.
var engineClasses = map[string]EngineClass{
	"Kestrel O-20": EngineClassA,
	"Kestrel D-30": EngineClassB,
	"Kestrel D-45": EngineClassC,
	"Kestrel D-60": EngineClassC,
	"Kestrel D-90": EngineClassC,
	"Kestrel D-120": EngineClassC,
	"Kestrel D-160": EngineClassD,
	"Kestrel D-220": EngineClassD,
	"Marlow D-20": EngineClassC,
	"Marlow O-30": EngineClassA,
	"Marlow D-45": EngineClassB,
	"Marlow D-60": EngineClassC,
	"Marlow D-90": EngineClassC,
	"Marlow D-120": EngineClassC,
	"Marlow D-160": EngineClassD,
	"Marlow D-220": EngineClassD,
	"Brent D-20": EngineClassB,
	"Brent D-30": EngineClassC,
	"Brent D-45": EngineClassC,
	"Brent D-60": EngineClassC,
	"Brent D-90": EngineClassC,
	"Brent D-120": EngineClassC,
	"Brent D-160": EngineClassD,
	"Brent D-220": EngineClassD,
	"Tarn O-20": EngineClassA,
	"Tarn D-30": EngineClassB,
	"Tarn D-45": EngineClassC,
	"Tarn D-60": EngineClassC,
	"Tarn D-90": EngineClassC,
	"Tarn D-120": EngineClassC,
	"Tarn D-160": EngineClassD,
	"Tarn D-220": EngineClassD,
	"Oyster D-20": EngineClassC,
	"Oyster O-30": EngineClassA,
	"Oyster D-45": EngineClassB,
	"Oyster D-60": EngineClassC,
	"Oyster D-90": EngineClassC,
	"Oyster D-120": EngineClassC,
	"Oyster D-160": EngineClassD,
	"Oyster D-220": EngineClassD,
	"Fulmar D-20": EngineClassB,
	"Fulmar D-30": EngineClassC,
	"Fulmar D-45": EngineClassC,
	"Fulmar D-60": EngineClassC,
	"Fulmar D-90": EngineClassC,
	"Fulmar D-120": EngineClassC,
	"Fulmar D-160": EngineClassD,
	"Fulmar D-220": EngineClassD,
	"Heron O-20": EngineClassA,
	"Heron D-30": EngineClassB,
	"Heron D-45": EngineClassC,
	"Heron D-60": EngineClassC,
	"Heron D-90": EngineClassC,
	"Heron D-120": EngineClassC,
	"Heron D-160": EngineClassD,
	"Heron D-220": EngineClassD,
}

// ClassOf returns the Schedule 2 class of an engine model.
func ClassOf(model string) (EngineClass, error) {
	c, ok := engineClasses[model]
	if !ok {
		return "", fmt.Errorf("ccvr: engine model %q has no Schedule 2 class", model)
	}
	return c, nil
}

// UnrecordedInput is an input a rule depends on that nothing in the tool can
// capture today. Each one is reported as an open finding until a command exists
// that records it; the rule then reads the recorded value instead.
type UnrecordedInput struct {
	Key     string
	Label   string
	Needs   string
	Assumes string
	Cites   []Citation
}

// Citation points at a clause of the regulation and carries the text the
// rule relies on.
type Citation struct {
	Clause string
	Text   string
}

// EngineHours is the rule behind reg. 11(1): hours run since the last
// overhaul, per engine, against the Schedule 2 threshold.
type EngineHours struct {
	Engine        string
	Class         EngineClass
	SinceOverhaul float64
	Threshold     float64
}

// Over reports whether the engine is past its threshold.
func (h EngineHours) Over() bool { return h.SinceOverhaul > h.Threshold }

func hoursTowardsOverhaul(voyages []VoyageHours, overhauls []Overhaul) ([]EngineHours, error) {
	last := map[string]time.Time{}
	for _, o := range overhauls {
		if o.At.After(last[o.Engine]) {
			last[o.Engine] = o.At
		}
	}
	byEngine := map[string]*EngineHours{}
	for _, v := range voyages {
		for engine, hours := range v.Hours {
			h, ok := byEngine[engine]
			if !ok {
				class, err := ClassOf(v.Models[engine])
				if err != nil {
					return nil, err
				}
				h = &EngineHours{Engine: engine, Class: class, Threshold: overhaulThreshold[class]}
				byEngine[engine] = h
			}
			if v.Start.After(last[engine]) {
				h.SinceOverhaul += hours
			}
		}
	}
	out := make([]EngineHours, 0, len(byEngine))
	for _, h := range byEngine {
		out = append(out, *h)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Engine < out[j].Engine })
	return out, nil
}

// extendedInterval is the survey interval in reg. 17(1): two seasons.
const extendedInterval = 2

// surveyDue applies reg. 17(1) and (2): the next survey is due two seasons
// after the date on the last survey certificate.
func surveyDue(lastSurvey time.Time) time.Time {
	return lastSurvey.AddDate(extendedInterval, 0, 0)
}

// LevyBand is a band in Schedule 3.
type LevyBand int

// levyBand applies reg. 25 and Schedule 3 to the season's engine hours.
func levyBand(seasonHours float64) LevyBand {
	switch {
	case seasonHours <= 200:
		return 1
	case seasonHours <= 600:
		return 2
	default:
		return 3
	}
}

// clauses renders the citations of a finding for its text.
func clauses(cs []Citation) string {
	names := make([]string, len(cs))
	for i, c := range cs {
		names[i] = c.Clause
	}
	return strings.Join(names, ", ")
}

// quotes holds the text each rule relies on, keyed by clause, so a test can
// check it against docs/regulation/ccvr-2021/text.md (ADR-0005).
var quotes = map[string]string{
	"CCVR 2021 reg. 1(1)": "These Regulations may be cited as the Coastal Craft (Voyage and Returns) Regulations 2021 (CCVR 2021).",
	"CCVR 2021 reg. 1(4)": "Where the vessel is out of commission for part of the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 1(8)": "Where the Authority has granted an exemption, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 1(12)": "Where a record is kept in a language other than English, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 1(16)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 2(23)": "Where a vessel has more than one engine, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 3(1)": "The coastal zone is the sea within twelve nautical miles of the coast, together with the harbours and estuaries that open onto it.",
	"CCVR 2021 reg. 3(5)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 3(9)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 3(13)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 3(17)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 3(20)": "Where the operator changes during the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 5(1)": "The operator of a vessel is the person who has its day-to-day management, whether or not that person owns it.",
	"CCVR 2021 reg. 5(5)": "Where a record required by this regulation is lost or destroyed, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 5(9)": "Where the operator is a company, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 5(13)": "Where a surveyor finds that a record is incomplete, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 5(17)": "Where a record required by this regulation is lost or destroyed, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 5(20)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 5(24)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 6(2)": "The return states the engine hours for the season, the hours towards overhaul for each engine, the survey due date and the levy band.",
	"CCVR 2021 reg. 6(6)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 6(10)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 6(13)": "Where the Authority has granted an exemption, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 6(17)": "Where a record is kept in a language other than English, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 6(21)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 6(25)": "Where the Authority has granted an exemption, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 6(29)": "Where a record is kept in a language other than English, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 7(3)": "Where a surveyor finds that a record is incomplete, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 7(7)": "Where a record required by this regulation is lost or destroyed, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 7(11)": "Where the operator is a company, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 7(15)": "Where a surveyor finds that a record is incomplete, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 7(19)": "Where a record required by this regulation is lost or destroyed, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 7(23)": "Where the operator is a company, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 7(26)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 8(3)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 8(7)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 8(11)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 8(15)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 8(18)": "Where a vessel is sold during the season, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 8(22)": "Where the vessel is chartered without crew, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 8(26)": "Where two records conflict, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 9(3)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 9(7)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 9(11)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 9(14)": "Where the operator changes during the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 9(18)": "Where a vessel has more than one engine, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 9(22)": "Where the operator cannot produce a record within 28 days of a request, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 9(26)": "Where the operator changes during the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 10(2)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 10(5)": "Where the Authority has granted an exemption, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 10(9)": "Where a record is kept in a language other than English, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 10(13)": "Where the vessel is out of commission for part of the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 10(17)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 10(21)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 10(24)": "Where the vessel is chartered without crew, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 11(3)": "Hours credited under paragraph (2) shall be evidenced by the engine log or, where the log is lost, by a statement of the operator.",
	"CCVR 2021 reg. 12(3)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 12(7)": "Where the Authority has granted an exemption, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 12(11)": "Where a record is kept in a language other than English, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 12(15)": "Where the vessel is out of commission for part of the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 12(18)": "Where two records conflict, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 12(22)": "Where a vessel is sold during the season, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 12(26)": "Where the vessel is chartered without crew, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 12(30)": "Where two records conflict, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 12(34)": "Where a vessel is sold during the season, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 12(37)": "Where the operator cannot produce a record within 28 days of a request, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 12(41)": "Where the operator changes during the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 12(45)": "Where a vessel has more than one engine, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 13(1)": "Where a record required by this regulation is lost or destroyed, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 13(5)": "Where the operator is a company, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 13(9)": "Where a surveyor finds that a record is incomplete, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 13(12)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 13(16)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 13(20)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 13(24)": "Where a record is kept in a language other than English, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 13(28)": "Where the vessel is out of commission for part of the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 13(31)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 13(35)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 13(39)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 13(43)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(3)": "Where the operator cannot produce a record within 28 days of a request, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 14(6)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(10)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(14)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(18)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(22)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(26)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 14(29)": "Where the vessel is out of commission for part of the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 14(33)": "Where the Authority has granted an exemption, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 14(37)": "Where a record is kept in a language other than English, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 14(41)": "Where the vessel is out of commission for part of the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 14(45)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 15(3)": "Where a record required by this regulation is lost or destroyed, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 15(7)": "Where the operator is a company, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 15(11)": "Where a surveyor finds that a record is incomplete, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 15(15)": "Where a record required by this regulation is lost or destroyed, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 15(19)": "Where the operator is a company, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 15(23)": "Where a surveyor finds that a record is incomplete, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 15(26)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 15(30)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 15(34)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 15(38)": "Where a record is kept in a language other than English, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 15(42)": "Where the vessel is out of commission for part of the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 15(45)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(1)": "Where a vessel has more than one engine, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 16(5)": "Where the operator cannot produce a record within 28 days of a request, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 16(9)": "Where the operator changes during the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 16(13)": "Where a vessel has more than one engine, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 16(16)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(20)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(24)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(28)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(32)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(36)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 16(39)": "Where a record is kept in a language other than English, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 16(43)": "Where the vessel is out of commission for part of the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 17(3)": "The extended survey interval in paragraph (1) does not apply to a vessel carrying passengers for hire or reward.",
	"CCVR 2021 reg. 18(2)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 18(6)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 18(9)": "Where the vessel is out of commission for part of the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 18(13)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 18(17)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 18(21)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 18(25)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 18(29)": "Where a record is kept in a language other than English, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 18(32)": "Where the vessel is chartered without crew, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 18(36)": "Where two records conflict, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 18(40)": "Where a vessel is sold during the season, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 18(44)": "Where the vessel is chartered without crew, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 19(3)": "Where a surveyor finds that a record is incomplete, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 19(6)": "Where a record is kept in a language other than English, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 19(10)": "Where the vessel is out of commission for part of the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 19(14)": "Where the Authority has granted an exemption, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 19(18)": "Where a record is kept in a language other than English, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 19(22)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 19(25)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 19(29)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 19(33)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 19(37)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(1)": "Where the operator changes during the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 20(5)": "Where a vessel has more than one engine, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 20(8)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(12)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(16)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(20)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(24)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 20(27)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 20(31)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 20(35)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 20(39)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 21(2)": "Where the operator changes during the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 21(6)": "Where a vessel has more than one engine, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 21(9)": "Where a record required by this regulation is lost or destroyed, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 21(13)": "Where the operator is a company, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 21(17)": "Where a surveyor finds that a record is incomplete, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 21(21)": "Where a record required by this regulation is lost or destroyed, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 21(25)": "Where the operator is a company, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 21(28)": "Where the Authority has granted an exemption, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 21(32)": "Where a record is kept in a language other than English, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 21(36)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 21(40)": "Where the Authority has granted an exemption, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 21(44)": "Where a record is kept in a language other than English, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 22(1)": "Hull work on a vessel to which these Regulations apply shall be carried out at a yard on the approved list published under regulation 31, save where paragraph (4) or (6) applies.",
	"CCVR 2021 reg. 22(6)": "Hull work carried out at a yard not on the approved list counts towards the survey interval only in the proportion that hours worked at approved yards bear to the total hours of hull work in the season. The operator shall hold, for work under this paragraph, the yard's invoice or a statement from the yard of the hours worked.",
	"CCVR 2021 reg. 23(3)": "Where the operator is a company, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 23(7)": "Where a surveyor finds that a record is incomplete, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 23(11)": "Where a record required by this regulation is lost or destroyed, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 23(15)": "Where the operator is a company, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 24(2)": "Where the vessel is chartered without crew, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 24(6)": "Where two records conflict, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 24(10)": "Where a vessel is sold during the season, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 24(14)": "Where the vessel is chartered without crew, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 24(18)": "Where two records conflict, the operator shall notify the Authority within 28 days.",
	"CCVR 2021 reg. 25(2)": "Where the operator cannot produce a record within 28 days of a request, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 25(6)": "Where the operator changes during the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 25(10)": "Where a vessel has more than one engine, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 25(14)": "Where the operator cannot produce a record within 28 days of a request, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 25(18)": "Where the operator changes during the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 26(4)": "Where the vessel is chartered without crew, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 26(7)": "Where the operator changes during the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 26(11)": "Where a vessel has more than one engine, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 26(15)": "Where the operator cannot produce a record within 28 days of a request, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 27(2)": "Where a record is kept in a language other than English, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 27(6)": "Where the vessel is out of commission for part of the season, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 27(9)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 27(13)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 27(17)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 28(3)": "Where a record is kept in a language other than English, the operator shall keep a translation with the record.",
	"CCVR 2021 reg. 28(7)": "Where the vessel is out of commission for part of the season, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 28(10)": "Where two records conflict, the requirement applies to each engine separately.",
	"CCVR 2021 reg. 28(14)": "Where a vessel is sold during the season, the duty falls on each person who was the operator during the season.",
	"CCVR 2021 reg. 28(18)": "Where the vessel is chartered without crew, the later record prevails unless the operator shows otherwise.",
	"CCVR 2021 reg. 29(4)": "Where a record is kept in a language other than English, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 29(8)": "Where the vessel is out of commission for part of the season, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 29(12)": "Where the Authority has granted an exemption, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 29(15)": "A record kept for the purposes of this regulation shall state:",
	"CCVR 2021 reg. 30(1)": "Where the Authority has granted an exemption, the operator shall make a signed statement of what the record contained.",
	"CCVR 2021 reg. 30(5)": "Where a record is kept in a language other than English, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 30(9)": "Where the vessel is out of commission for part of the season, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 30(13)": "Where the Authority has granted an exemption, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 30(16)": "Where a vessel is sold during the season, a surveyor may treat the record as not kept.",
	"CCVR 2021 reg. 31(2)": "Where the Authority has granted an exemption, the requirement applies to the part of the season in which the vessel was in commission.",
	"CCVR 2021 reg. 31(6)": "Where a record is kept in a language other than English, the exemption is stated on the operating return.",
	"CCVR 2021 reg. 31(10)": "Where the vessel is out of commission for part of the season, the Authority may accept a statement of the operator instead.",
	"CCVR 2021 reg. 31(14)": "Where the Authority has granted an exemption, the operator shall keep a translation with the record.",
}

// registry lists the clauses that offer the operator a choice, in the order
// `tidewright clauses` prints them.
func registry() []Clause {
	return []Clause{
		// reg. 18(2): a survey may be deferred once, by up to three months.
		{Key: ClauseSurveyDeferral,
			Label:   "Survey deferral",
			Options: []Option{{Key: "none"}, {Key: "three-months"}},
		},
		// reg. 19(1): a season laid up counts or does not count towards the interval.
		{Key: ClauseLaidUpSeason,
			Label:   "Laid-up season",
			Options: []Option{{Key: "counts"}, {Key: "excluded"}},
		},
		// reg. 26(3): levy paid at once or in two instalments.
		{Key: ClauseLevyInstalments,
			Label:   "Levy instalments",
			Options: []Option{{Key: "single"}, {Key: "two"}},
		},
		// reg. 12(2): overhaul records kept per engine or per vessel.
		{Key: ClauseOverhaulRecord,
			Label:   "Overhaul record",
			Options: []Option{{Key: "per-engine"}, {Key: "per-vessel"}},
		},
		// reg. 28(1): a correction replaces the return or is filed as an addendum.
		{Key: ClauseCorrectionMethod,
			Label:   "Correction method",
			Options: []Option{{Key: "replace"}, {Key: "addendum"}},
		},
		// reg. 7(2): season starts on 1 January or on first commissioning.
		{Key: ClauseSeasonStart,
			Label:   "Season start",
			Options: []Option{{Key: "calendar"}, {Key: "commissioning"}},
		},
		// reg. 14(3): hull class by length or by displacement.
		{Key: ClauseHullClass,
			Label:   "Hull class basis",
			Options: []Option{{Key: "length"}, {Key: "displacement"}},
		},
		// reg. 20(2): a tender is its own vessel or part of its parent.
		{Key: ClauseTenderAsVessel,
			Label:   "Tender as a separate vessel",
			Options: []Option{{Key: "separate"}, {Key: "part-of-parent"}},
		},
	}
}

var unrecordedInputs = []UnrecordedInput{
	{Key: InputHomePort,
		Label:   "Home port",
		Needs:   "mooring link: home port of the vessel during the 2025 season",
		Assumes: "Assumes the vessel is home-ported inside the coastal zone all season.",
		Cites: []Citation{{Clause: "CCVR 2021 reg. 4(1)",
			Text: "These Regulations apply to a vessel whose home port is within the coastal zone for the whole or any part of the season."}},
	},
	{Key: InputForHire,
		Label:   "Operated for hire",
		Needs:   "operator fact or vessel attribute: whether the vessel carries passengers for hire or reward",
		Assumes: "Assumes no paying passengers, so the extended survey interval is used.",
		Cites: []Citation{{Clause: "CCVR 2021 reg. 17(3)",
			Text: "The extended survey interval in paragraph (1) does not apply to a vessel carrying passengers for hire or reward."}},
	},
	{Key: InputOffListYardHours,
		Label:   "Off-list yard work",
		Needs:   "hand entry: hours of hull work carried out at yards not on the approved list during the season",
		Assumes: "Assumes every hour of hull work was done at an approved yard.",
		Cites: []Citation{{Clause: "CCVR 2021 reg. 22(6)",
			Text: "Hull work carried out at a yard not on the approved list counts towards the survey interval only in the proportion that hours worked at approved yards bear to the total hours of hull work in the season."}},
	},
	{Key: InputCarriedEngineHours,
		Label:   "Carried engine hours",
		Needs:   "carried balance: engine hours run in the five preceding seasons, before the period the imported logs cover",
		Assumes: "Assumes no hours carried in: hours towards overhaul start at zero on the first imported voyage.",
		Cites: []Citation{{Clause: "CCVR 2021 reg. 11(2)"}, {Clause: "CCVR 2021 reg. 11(5)",
			Text: "Engine hours run in any of the five preceding seasons may be credited against the overhaul threshold: in each season by no more than 50 per cent of those hours, or once by an amount not exceeding 2,000 hours."}},
	},
}

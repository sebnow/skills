# Changelog

## Unreleased

- A room can set `max_minutes`; a longer booking of that room is refused
  with 422. Rooms without it have no cap.
- A room can set `members_only`; only users listed under `members` can book
  it, and others are refused with 403.
- `POST /bookings/check` takes the same body as `POST /bookings` and answers
  with the status and error the booking would get, or 200 if it would be
  accepted, without booking anything. It does not hold the room.

## 0.3.0

- `DELETE /bookings/{id}` lets the member who made a booking cancel it.

## 0.2.0

- `GET /rooms/{id}/bookings` lists a room's bookings that have not ended,
  earliest first.

## 0.1.0

- Book a room with `POST /bookings`. A booking that overlaps another booking
  of the same room is refused with 409.
- Rooms come from `deskbook.json`; see `deskbook.example.json`.
- deskbook runs behind the office SSO proxy, which sets `X-Forwarded-User` on
  every request; requests without it are refused with 401.
- Bookings are kept in memory.

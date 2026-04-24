# Missing Backend Features

## Invitation System (Designed but not implemented)

The `Invitation` model and `InvitationRepository` port exist in core/domain and core/ports/outbound respectively, but no HTTP handlers/routes expose invitation management.

### What's needed:
- `POST /invitations` — Create invitation (Admin/Member → Guest)
- `GET /invitations` — List invitations (Admin sees all, Members see their own)
- `PATCH /invitations/:id` — Accept/Decline invitation
- `GET /invitations/:id` — Get invitation details

### Register flow update:
Currently: `/auth/register` uses env-based `ADMIN_INVITE_CODE`
Should be: `/auth/register` uses an `Invitation` record with status=pending

### Invite management UI:
- Admin: "Manage Invitations" panel — create invite, view pending/accepted/declined
- Member: "Invite Guest" — creates pending invitation
- Guest: sees pending invitation to accept/decline

---

## League play dates

`POST /leagues/:id/play-dates` is listed in API.md but not implemented in router.go.

---

## League PATCH endpoint

`PATCH /leagues/:id` is listed in API.md but not in router.go.

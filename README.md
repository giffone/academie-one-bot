# Academie One Bot

The program is designed for building security pass system.

- Creating an administrator
- Creation of security personnel
- Creating invitations:
  - for students (with course completion marking)
  - for guests
- Student self-registration (course enrollment) by invitation
- Self-registration by guest (visiting the building) by invitation
- QR scanning at the entrance
- Guard post (notification of entrants to the building)

## Getting started

To run locally:

```bash
# env
- export DATABASE_URL="postgres://admin:password@localhost:5432/db_shema?search_path=telegram"
- export BOT_ID="000000000" # telegram_id of bot
- export BOT_TOKEN="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" # telegram bot token
- export WEB_URL="https://webapp.app/" # address of web application (frontend)
- export QR_TEXT="https://t.me/academie_one_bot" # text in qr code
- export REQ_LOG=false # logging all incoming requests (not realised)

# run application
make run
```

Frontend realised: [Ac Bot netlify](https://github.com/giffone/ac_bot_nav_netlify)

To run prod - use **CI/CD variables**

Gitlab build with **.gitlab-ci.yml**

## Application

### Student

![student](./z_readme/student.jpg)

### Guest

![guest](./z_readme/guest.jpg)

### Guard post

![guard post](./z_readme/post.jpg)

### Admin

![guard post](./z_readme/admin.jpg)
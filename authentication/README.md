# Authentication to Hoverfly

Hoverfly uses a combination of basic auth and JWT (JSON Web Tokens) to authenticate users

## Usage

Add new user:

    ./hoverfly -v -add -username hfadmin -password hfadminpass 

Getting token:

    curl -H "Content-Type application/json" -X POST -d '{"Username": "hoverfly", "Password": "testing"}' http://localhost:8888/token-auth

Using token:

    curl -H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTY0MDA2NzAsImlhdCI6MTQ1NjE0MTQ3MCwic3ViIjoiIn0.Z95DlJFP1nBTRQCpK_AkNYJvUqJpYLGijxttoqyAaf5hbfx1HML_I5uTZqWCK5oayVITih7P-zULA9DFCDQOTiLwRLIuJpsg9-9ApArSWQg-8JnrRk1IR1OjcYyPslQ5Zj1lat6AehYsmHb7A6-AWqKXF0_1XTz7lAZ2F_YExj9LrkSpkWqo4Qy58IltfjZXxwOOPdp7y2TmcjM8mpuF5sDxD9uCh74ahSsEnZxxVpIgHJJb9gQ3ZjYTPH8-h-yavINQ6ctl0Za-oaqG7tRdR3M5UH-eaBFDnEFzn7XdYwxyisiptXdULewt_KghpxzlloPMZsDRZsGs8VH6XxHDHg" http://localhost:8888/records

# A deep dive in distributed systems (in progress)

## Dev setup almost complete for auth and mailing'.'

    ## About auth':'
    ### auth should
    1. Authorize users
    2. Store sensitive user Data securely
    3. Store sensitive session Data securely
    4. Have a fallback if user forgotten their passwords
    5. Have some access roles
    6. Start the session
    7. Refresh the session
    8. End the session if the user want to end the session
    9. Some session fallbacks with expiry
    10. Interact with the mailer if there some access to verify

    ## About mailer_server':'
    ### mailer should 
    1. Receive some email data
    2. Should use the right template for the email data
    3. Send the email securely

    ## Done':'
    ### Component lib setup in auth
    ### Component lib tests in auth
    ### Component lib setup in mailer_server

    ## TODO':'
    ### Component lib setup in mailer_server
    ### Finish the configs for Docker-Compose
    ### Docker-Compose setup for the complete app
    ### Migrate all together
    ### Make a test environment for the whole infrastructure 

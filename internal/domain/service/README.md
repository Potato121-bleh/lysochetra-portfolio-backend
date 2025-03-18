# Instruction for the Constructing & usage of Services

# System infrastructure

-   Service are built for code smell purposes & provide a dynamic usability. which make all method of each service serve a purpose in Handler where the handler doesn't need to worried about the operation of database interaction or any type of business logic.

-   User_service

    This service will play a role to handle database interaction logic. They require handler to provide transaction (pgx.Tx) and additional args to met the requirement of services.

    Why service require us to make our own transaction and provide to the service?

    -   Because when handler is the one managing the whole transaction can make the performace more efficiency. As the service done their job, it the handler job to decide to commit or not. This feature is not built for continuous transaction. But it was built for handler to commit during the http response back to client.
        BUT IF the handler want to perform a serial transaction where it require to execution query a lot of time in a single transaction they can create new method to service to handle it.

    THis design pattern are fit for the scenario like this:

    -   When we don't make handler to manage the transaction, let say the service done their job which now the transaction also commit, but as the handler about to do http response back to client, the network got cut out, which cause a response mostly to be 500 status which is an error. SO this cause based on user POV they saw it error so they could login again while the transaction already complete without letting user know it. which can cause more duplicate later on future.

    -   So the solution is allow handler to manage the transaction where it commit AFTER the http response has sent.

    # System work flow:

    -   It checking for the tx (pgx.Tx) provided by handler. This tx will get checked if it nil (because handler want to do one-time-transaction) then the PrepTx will create new transaction which will get committed at the end of execution of service
    -   As the tx already clean, now the service call the repo to do the task and make sure it work well, if it failed it will rollback instantly and return error to handler
    -   As if the repo work well, service will call FinalizeTx to decide whether to commit transaction or not, by taking the provided tx from handler if it nil (meaning to do one-time-transaction) then we will commit the transaction if not then do nothing.
        Retrun :
        True: if everything work as expected
        False: commit or rollback is failed

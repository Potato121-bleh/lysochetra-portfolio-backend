package db

In this package it have all of tool you need for database including: - Interface for Database, DatabaseTx & Rows by pgx.Rows - Adapter for both Database & DatabaseTx Interface

As our DDD pattern design we design our system to met the requirement for both use case for database interact & testibility.

As you already know the system doesn't recongise other specific db tool such as pgxpool or pgx.Row or pgx.Tx. Because they use their own interface as mentioned above.

Because of that we provided with an adapter for both use case

-   so later on if we have an issue to use pgxpool.Pool for db we can use an adapter in database.go you can use:
    -   NewPgxDBAdapter() just passed in the db \*pgxpool.Pool the return db is good to go
    -   NewPgxTxAdapter() just passed in the tx pgx.Tx the return tx is good to go

⚠️⚠️⚠️ NOTE: As the whole system using interface for any access to db, tx or row which cause a limitation of features.
SO if you planned on using other feature beside pre-built interface you can add or unlock more feature by following:

-   Add new method on interface where you want to use.
-   Add those feature to all related mock  
     (Because mock also struct so they need to implement interface)

⭐ Example:
You want to use Reset() on pgxpool.Pool so you can go head to add new method in interface of Database: Reset() ...
Then you can go to MockDB and create new method inside of it (Reset() in our case).

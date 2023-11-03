# Test Scenarios

This folder includes tests added for seat-reservation app.

## Double Reservation Prevention

This section will describe which scenarios are covered in tests for double reservation prevention.
After brainstorming, following scenarios covered:

- **Scenario 1**: One user send same seat reservation and they are proccessed concurrently. In this scenario, expected behavior is to return success
- **Scenario 2**: Two or more users wants to reserve same seat and they are executed concurrently. In this scenario, first executed will succeed and second will fail
- **Scenario 3**: First user cancels a reservation and at the same time second user keeps on trying to reserve the same seat. Second user should succeed reserving after first user reservation successfully canceled

There are other possible scenarios to also consider that are hard to test:

- **Scenario 4**: Server disconnects right after creating transaction and performing part of the create reservation operation. In this scenario, transaction will be rolled back in postgres direct after server disconnection
- **Scenario 5**: Server creates DB transaction and goes into infinite loop before finishing transaction. In this scenario, postgresDB will rollback after `idle_in_transaction_session_timeout` period

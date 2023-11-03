# Test Scenarios

This folder includes tests added for seat-reservation app.

## Double Reservation Prevention

This section will describe which scenarios are covered in tests for double reservation prevention.
The way we solve double reservation prevention is by creating transaction before performing reservation operation.
Following would be the steps:

1. Begin transaction
2. Check if user has reservation for given event, If yes, rollback and return the reservation. Else, move to step (3)
3. Query seat. If seat is reserved, rollback. Else move to step (4)
4. Create reservation. Note that `seatID` is unique key in reservation table.
5. If creating reservation success, then update seat row to reserved since reservation is success(There will be always 1 reserver of the seat when reserved).
6. Commit transaction.

Similar idea should be applied for canceling the reservation. For this, we need to perform step (4) - (5) reversely: (5) - (4), as follow:

1. Begin transaction
2. Get reservation for given event id and user id. If no reservation, rollback and return. Else move to step (3).
3. Get the seat for given reservation and update seat to available.
4. Delete the reservation.
5. Commit transaction.

After brainstorming, following scenarios covered:

- **Scenario 1**: One user send same seat reservation and they are proccessed concurrently. In this scenario, expected behavior is to return success
- **Scenario 2**: Two or more users wants to reserve same seat and they are executed concurrently. In this scenario, first executed will succeed and second will fail
- **Scenario 3**: User keeps on trying to cancel the item. Only one cancel should succeed and all others should fail.
- **Scenario 4**: First user keeps on trying to cancel a reservation and at the same time second user keeps on trying to reserve the same seat. Second user should succeed reserving after first user reservation successfully canceled and final state of the reservation should be reserved not `canceled`.

There are other possible scenarios to also consider that are hard to test:

- **Scenario 5**: Server disconnects right after creating transaction and performing part of the create reservation operation. In this scenario, transaction will be rolled back in postgres direct after server disconnection
- **Scenario 6**: Server creates DB transaction and goes into infinite loop before finishing transaction. In this scenario, postgresDB will rollback after `idle_in_transaction_session_timeout` period
- **Scenario 7**: It is possible that one user is canceling and other trying to reserve the seat. If they are done concurrently, there is a  chance that one is in the middle of cancel fails before finishing the all cancelation steps giving way out to other user to reserve it. For example, we could delete reservation of the user and not updated all tables to correct state and other user can create reservation.

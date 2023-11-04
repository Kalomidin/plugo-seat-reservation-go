# Test Scenarios

This folder includes tests added for seat-reservation app.

## Double Reservation Prevention

### Method

Our approach involves initiating a database transaction before executing create or cancel reservation operations.

**Note**: In this solution, I do not assume that same one user will keep on trying to cancel and create reservation in parallel.

### Create Reservation Sequence

Following would be the steps:

1. Begin the transaction.
2. Query the database to determine if the user already has a reservation for the given event. If a reservation exists, rollback the transaction and return the existing reservation. If not, proceed to the next step.
3. Retrieve the status of the seat from the database.
4. If the seat is already reserved, rollback the transaction. If available, proceed to the next step.
5. Create a new reservation. Note: the seatID is a unique key in the reservation table.
6. If the reservation is successful, update the seat record to indicate it's reserved.
7. Commit the transaction.

### Cancel Reservation Sequence

1. Begin the transaction.
2. Retrieve the reservation using the provided event id and user id. If no reservation is found, rollback the transaction and return an error. Otherwise, proceed to the next step.
3. Retrieve the associated seat for the reservation.
4. If the seat is available, rollback the transaction, indicating a previous cancellation. If reserved, proceed to the next step.
5. Delete the reservation.
6. If the reservation is successful, update the seat record to indicate it's available.
7. Commit the transaction

## Scenarios

### Tested Scenarios

After brainstorming, following scenarios covered:

- **Scenario 1**: One user send same seat reservation and they are proccessed concurrently. In this scenario, expected behavior is to return success
  - Above will hanlde this since `seatId` is unique index in `reservation` table.
- **Scenario 2**: Two or more users wants to reserve same seat and they are executed concurrently. In this scenario, first executed will succeed and second will fail
  - Since we have unique index in `reservation` table for `seatId`, this will be handled as expected
- **Scenario 3**: User keeps on trying to cancel the item. Only one cancel should succeed and all others should fail.
  - A user can delete only once a reservation. Since multiple deletion will result to failure, it will behave as expected
- **Scenario 4**: First user keeps on trying to cancel a reservation and at the same time second user keeps on trying to reserve the same seat. Second user should succeed reserving after first user reservation successfully canceled and final state of the reservation should be reserved not `available`.
  - We have a strategy to update seat status after reservation is created/deleted. This implies that, if there exist a reservation record for seat, we can delete or create reservation only after creating or deleting reservation is **finished and seat is updated**,respectively. Thus, second user can not create reservation unless first user finishes the cancel operation and consequent cancel tasks from first user will just be scenario (2) and second user consequent tasks would behave as scenario (1).

### Hard-to-test Scenarios

There are other possible scenarios to also consider that are hard to test:

- **Scenario 5**: Server disconnects right after creating transaction and performing part of the create reservation operation. In this scenario, transaction will be rolled back in postgres direct after server disconnection but some other DB can fail.
- **Scenario 6**: The server acquires a significant lock on the database (e.g., table-level lock) and subsequently enters an infinite loop in the business logic, preventing the transaction from completing. In the case of PostgreSQL, the transaction will eventually be rolled back after reaching the idle_in_transaction_session_timeout period. However, for some other databases, this behavior could cause other processes to hang indefinitely, awaiting the release of the lock.
- **Scenario 7**: It is possible that one user is canceling and other trying to reserve the seat. If they are done concurrently, there is a  chance that one is in the middle of cancel fails before finishing the all cancelation steps giving way out to other user to reserve it. For example, we could delete reservation of the user and not updated all tables to correct state and other user can create reservation.

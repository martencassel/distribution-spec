# Purpose: Requesting a token

## Functions:

1. Verifies a users identity

Questions answered: Is the user identified ?
- Accepts Basic Auth (username/password or PAT).
- Validates credentials against the registry servers identity backend.
- Rejects invalid or missing credentials.

2. Validates the requested scope

Question answered: Does the user have permissions for the requested operations ?
- Parses the scope parameter:
	repository:samalba/my-app:pull,push
- Checks whether the authenticated user has rights to perform the requested action on the repository.
- Denies the request if permissions do not match.

3. Issues a bearer token

Purpose: Provides a short-time-lived security token

- Generates an opaque Bearer token
- Encodes allowed actions (pull/push) in the token’s access section.
- Returns metadata such as expires_in and issued_at.
This token is later used by the registry to authorize operations.

4. Enables the registry to allow or deny operations
Purpose: Let the registryenforce access control.

- The client includes the token in: 	Authorization: Bearer <token>
- The registry validates the token and checks whether the requested operation matches the token’s granted scope.
- The registry either:
	allows the operation (pull/push), or
	denies it if the token is invalid or insufficient.

# Flow

1. Client makes a pull request
2. Registry responds with 401 Unauthorized
3. Client fetches a token
4. Auth service processes the request
5. Client retries the oriinal request
6. Registry validates the token


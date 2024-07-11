Library Management SystemRequirements:

Books Management:
Each book has a title, author, and a unique ISBN.
Books can be added, removed, and listed.

User Management:
Each user has a name and a unique ID.
Users can be added, removed, and listed.

Borrowing and Returning Books:
Users can borrow books if they are available.
Users can return books.
More than one user cannot borrow a book at a time.

Concurrency:
Ensure thread safety when multiple users are borrowing and returning books simultaneously.

Error Handling:
Proper error messages should be returned when operations fail (e.g., borrowing a non-existent book, returning a book not borrowed).

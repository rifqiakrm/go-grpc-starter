## Making Changes

- Clone the project

- If there are some interfaces you want to create some mocks for them, you can use `mockgen`. We provide a simple command to generate necessary mock.
  Follow [How to Create Mock](CREATE_MOCK.md).

- If there are needs to migrate the database, follow [Database Migration](DATABASE_MIGRATION.md).

- If you have any SQL query in your changes, make sure your query includes the module's schema. For example, instead of using this query:

    ```
    SELECT * FROM questions_packets;
    ```

  you should use:

    ```
    SELECT * FROM assessment.questions_packets;
    ```

  Apply that to all of your SQL queries, not just SELECT, but also INSERT, UPDATE, DELETE, etc.

- Make sure you format/beautify the code by running

    ```
    $ make pretty
    ```

- As a reminder, always run the command above before add and commit changes.
  That command will be run in CI Pipeline to verify the format.

- If your changes require any document to be updated, please update the relevant documents.

- Add, commit, and push the changes to repository

    ```
    $ git add .
    $ git commit -s -m "your meaningful commit message"
    $ git push origin <your-meaningful-branch>
    ```

  For writing commit message, please use [conventionalcommits](https://www.conventionalcommits.org/en/v1.0.0/) as a reference.

- Create a Merge Request (MR). In your MR's description, please explain the goal of the MR and its changes.

- Ask the other contributors to review.

- Once your MR is approved and its pipeline status is green, merge your MR.
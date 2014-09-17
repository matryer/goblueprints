To start the envionrment in which `counter` and `twittervotes` will run:

Each in its own Terminal window/tab:

    mongod

    nsqlookupd

    nsqd --lookupd-tcp-address=localhost:4160

Then start `counter` and `twittervotes` in their own terminal windows too.

Open MondoDB shell:

    mongo

Create a new poll in the `ballots` database:

    use ballots
    db.polls.insert({title:"My poll",options:["one","two","three"]})

After a while, see the results by printing the polls:

    db.polls.find().pretty()

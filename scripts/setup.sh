docker exec -it mongo1 --eval "rs.initiate(
  {
    _id : 'myReplicaSet',
    members: [
      { _id : 0, host : "mongo1:30001" },
      { _id : 1, host : "mongo2:30002" },
      { _id : 2, host : "mongo3:30003" }
    ]
  }
)"

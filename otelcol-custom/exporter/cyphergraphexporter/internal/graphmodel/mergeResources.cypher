// param $resources is a list of newResourceParam structs with the fields 
UNWIND $resources AS resource
UNWIND resource.contained_in AS contained_node
MERGE (r:Resource {type: resource.label, id: resource.id})-[:IS_CONTAINED]->(c:Resource {type: contained_node.label, id: contained_node.id})
  ON CREATE
    SET
      r.created = timestamp(),
      r += resource.props
  ON MATCH
    SET
      r.lastSeen = timestamp()

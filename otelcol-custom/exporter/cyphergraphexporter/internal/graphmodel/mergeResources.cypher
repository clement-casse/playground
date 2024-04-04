// param $resources is a list of newResourceParam structs with the fields 
UNWIND $resources AS resource
MERGE (r:Resource {type: resource.label, id: resource.id})
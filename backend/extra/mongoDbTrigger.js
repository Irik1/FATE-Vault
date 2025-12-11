// Trigger function to update createdAt and updatedAt fields on MongoDB documents on insert and update operations of all collection in main db.

exports = function(changeEvent) {
  service = context.services.get("Cluster_name");
  const db = service.db("Db_name");
  const collection = db.collection(changeEvent.ns.coll);
  
  // Get the current timestamp
  const now = new Date();
  
  // Handle insert operations
  if (changeEvent.operationType === "insert") {
    const document = changeEvent.fullDocument;
    
    // Only update if createdAt or updatedAt are not already set
    const updateFields = {};
    
    if (!document.createdAt) {
      updateFields.createdAt = now;
    }
    
    if (!document.updatedAt) {
      updateFields.updatedAt = now;
    }
    
    // Only perform update if there are fields to update
    if (Object.keys(updateFields).length > 0) {
      collection.updateOne(
        { _id: document._id },
        { $set: updateFields }
      );
    }
  }
  
  // Handle update operations
  if (changeEvent.operationType === "update") {
    const documentId = changeEvent.documentKey._id;
    
    // Check if updatedAt is being explicitly set in the update
    const updateDescription = changeEvent.updateDescription || {};
    const updatedFields = updateDescription.updatedFields || {};
    
    // Only update updatedAt if it's not already being set in the update operation
    if (!updatedFields.hasOwnProperty("updatedAt")) {
      collection.updateOne(
        { _id: documentId },
        { $set: { updatedAt: now } }
      );
    }
  }
  
  // Handle replace operations (treat similar to update)
  if (changeEvent.operationType === "replace") {
    const document = changeEvent.fullDocument;
    
    if (document) {
      const updateFields = {};
      
      // Set createdAt only if it doesn't exist
      if (!document.createdAt) {
        updateFields.createdAt = now;
      }
      
      // Always update updatedAt on replace
      updateFields.updatedAt = now;
      
      collection.updateOne(
        { _id: document._id },
        { $set: updateFields }
      );
    }
  }
};
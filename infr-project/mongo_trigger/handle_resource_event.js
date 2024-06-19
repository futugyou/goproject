exports = async function (changeEvent) {
    const serviceName = "Cluster0";
    const dbName = "infr-project";
    const resourceChangeEventCollectionName = "resource_events";
    const resourceQueryCollectionName = "resources_query";

    // Initialize the ResourceQuery collection if not already done
    await initializeResourceQueryCollection();

    const resourceQueryCollection = context.services.get(serviceName).db(dbName).collection(resourceQueryCollectionName);

    // Check if the full document is available in the change event
    if (!changeEvent.fullDocument) {
        console.error("No fullDocument in changeEvent");
        return;
    }

    const resourceChangeEvent = changeEvent.fullDocument;
    const resourceId = resourceChangeEvent.id;
    const updateFields = {};

    if ('version' in resourceChangeEvent) updateFields.version = resourceChangeEvent.version;
    if ('type' in resourceChangeEvent) updateFields.type = resourceChangeEvent.type;
    if ('data' in resourceChangeEvent) updateFields.data = resourceChangeEvent.data;
    if ('name' in resourceChangeEvent) updateFields.name = resourceChangeEvent.name;
    if ('created_at' in resourceChangeEvent) updateFields.updated_at = resourceChangeEvent.created_at;
    if ('tags' in resourceChangeEvent) updateFields.tags = resourceChangeEvent.tags;
    updateFields.is_deleted = resourceChangeEvent.is_deleted ?? false;

    try {
        if (changeEvent.ns.coll === resourceChangeEventCollectionName) {
            // Find the current ResourceQuery document
            const currentResourceQuery = await resourceQueryCollection.findOne({ id: resourceId });

            if (currentResourceQuery) {
                // If the existing ResourceQuery document exists, update it with the new change event
                if (updateFields.is_deleted) {
                    await resourceQueryCollection.updateOne(
                        { id: resourceId },
                        { $set: { is_deleted: updateFields.is_deleted, updated_at: updateFields.updated_at } }
                    );
                } else {
                    await resourceQueryCollection.updateOne(
                        { id: resourceId },
                        { $set: updateFields }
                    );
                }
            } else {
                // If the existing ResourceQuery document does not exist, create a new one
                const newResourceQuery = {
                    id: resourceId,
                    ...updateFields,
                    created_at: updateFields.updated_at,
                    is_deleted: updateFields.is_deleted
                };

                await resourceQueryCollection.insertOne(newResourceQuery);
            }
        }
    } catch (err) {
        console.error("Error performing MongoDB write operation: ", err.message);
    }
};

async function initializeResourceQueryCollection() {
    const serviceName = "Cluster0";
    const dbName = "infr-project";
    const resourceChangeEventCollectionName = "resource_events";
    const resourceQueryCollectionName = "resources_query";
    const initializationFlagCollectionName = "InitializationFlags";
    const initializationFlagDocumentId = "resourceQueryInitialization";

    const resourceChangeEventCollection = context.services.get(serviceName).db(dbName).collection(resourceChangeEventCollectionName);
    const resourceQueryCollection = context.services.get(serviceName).db(dbName).collection(resourceQueryCollectionName);
    const initializationFlagCollection = context.services.get(serviceName).db(dbName).collection(initializationFlagCollectionName);

    // Check if the initialization has already been done
    const initializationFlag = await initializationFlagCollection.findOne({ _id: initializationFlagDocumentId });

    if (initializationFlag) {
        console.log("ResourceQuery collection has already been initialized.");
        return;
    }

    try {
        // Fetch all unique resource IDs from ResourceChangeEvent collection
        const resourceChangeEvents = await resourceChangeEventCollection.aggregate([
            {
                $group: {
                    _id: "$id",
                    events: { $push: "$$ROOT" }
                }
            },
            {
                $sort: { "events.version": 1 }
            }
        ]).toArray();

        for (let resourceChangeEvent of resourceChangeEvents) {
            const resourceId = resourceChangeEvent.id;
            const events = resourceChangeEvent.events;

            if (events.length > 0) {
                const latestEvent = events[events.length - 1];
                const firstEvent = events[0];
                const is_deleted = events.some(item => item.is_deleted === true);

                // Construct the resource query document
                const resourceQuery = { id: latestEvent.id };

                if ('name' in latestEvent) resourceQuery.name = latestEvent.name;
                if ('version' in latestEvent) resourceQuery.version = latestEvent.version;
                if ('type' in latestEvent) resourceQuery.type = latestEvent.type;
                if ('data' in latestEvent) resourceQuery.data = latestEvent.data;
                if ('tags' in latestEvent) resourceQuery.tags = latestEvent.tags;
                if ('created_at' in firstEvent) resourceQuery.created_at = firstEvent.created_at;
                if ('created_at' in latestEvent) resourceQuery.updated_at = latestEvent.created_at;
                resourceQuery.is_deleted = is_deleted;

                // Upsert the document in ResourceQuery collection
                await resourceQueryCollection.updateOne(
                    { id: resourceQuery.id },
                    { $set: resourceQuery },
                    { upsert: true }
                );
            }
        }

        // Set the initialization flag
        await initializationFlagCollection.insertOne({ _id: initializationFlagDocumentId, initialized: true });

        console.log("Initialization of ResourceQuery collection complete.");
    } catch (err) {
        console.error("Error initializing ResourceQuery collection: ", err.message);
    }
}

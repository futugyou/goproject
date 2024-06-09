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
    const version = resourceChangeEvent.version;
    const type = resourceChangeEvent.type;
    const data = resourceChangeEvent.data;
    const name = resourceChangeEvent.name;
    const created_at = resourceChangeEvent.created_at;
    const is_deleted = resourceChangeEvent.is_deleted ?? false;

    try {
        if (changeEvent.ns.coll === resourceChangeEventCollectionName) {
            // Find the current ResourceQuery document
            const currentResourceQuery = await resourceQueryCollection.findOne({ id: resourceId });

            if (currentResourceQuery) {
                // If the existing ResourceQuery document exists, update it with the new change event
                await resourceQueryCollection.updateOne(
                    { id: resourceId },
                    {
                        $set: { name: name, version: version, type: type, data: data, updated_at: created_at, is_deleted: is_deleted }
                    }
                );
            } else {
                // If the existing ResourceQuery document does not exist, create a new one
                const newResourceQuery = {
                    id: resourceId,
                    name: name,
                    version: version,
                    type: type,
                    data: data,
                    created_at: created_at,
                    is_deleted: is_deleted
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
                const is_deleted = events.some(item => item.is_deleted === true)
                // Construct the resource query document
                const resourceQuery = {
                    id: latestEvent.id,
                    name: latestEvent.name,
                    version: latestEvent.version,
                    type: latestEvent.type,
                    data: latestEvent.data,
                    created_at: firstEvent.created_at,
                    is_deleted: is_deleted
                };

                if (events.length > 1) {
                    resourceQuery.updated_at = latestEvent.created_at;
                }
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

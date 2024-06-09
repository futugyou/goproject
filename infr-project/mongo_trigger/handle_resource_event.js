exports = async function (changeEvent) {
    const serviceName = "Cluster0";
    const dbName = "infr-project";
    const resourceChangeEventCollectionName = "resource_events";
    const resourceQueryCollectionName = "resources_query";

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

    try {
        if (changeEvent.ns.coll === resourceChangeEventCollectionName) {
            // Find the current ResourceQuery document
            const currentResourceQuery = await resourceQueryCollection.findOne({ id: resourceId });

            if (currentResourceQuery) {
                // If the existing ResourceQuery document exists, update it with the new change event
                await resourceQueryCollection.updateOne(
                    { id: resourceId },
                    {
                        $set: { name: name, version: version, type: type, data: data, update_at: created_at }
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
                    created_at: created_at
                };

                await resourceQueryCollection.insertOne(newResourceQuery);
            }
        }
    } catch (err) {
        console.error("Error performing MongoDB write operation: ", err.message);
    }
};

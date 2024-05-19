const { MongoClient } = require('mongodb');

async function main() {
    const uri = "mongodb://localhost:27017";
    const client = new MongoClient(uri);
    try {
        await client.connect();
        const database = client.db('statesdb');
        const collection = database.collection('states');
        
        const states = [
            { name: "Alabama" },
            { name: "Alaska" },
            { name: "Arizona" },
            { name: "Arkansas" },
            { name: "California" },
            { name: "Colorado" },
            { name: "Connecticut" },
            { name: "Delaware" },
            { name: "Florida" },
            { name: "Georgia" },
            { name: "Hawaii" },
            { name: "Idaho" },
            { name: "Illinois" },
            { name: "Indiana" },
            { name: "Iowa" },
            { name: "Kansas" },
            { name: "Kentucky" },
            { name: "Louisiana" },
            { name: "Maine" },
            { name: "Maryland" },
            { name: "Massachusetts" },
            { name: "Michigan" },
            { name: "Minnesota" },
            { name: "Mississippi" },
            { name: "Missouri" },
            { name: "Montana" },
            { name: "Nebraska" },
            { name: "Nevada" },
            { name: "New Hampshire" },
            { name: "New Jersey" },
            { name: "New Mexico" },
            { name: "New York" },
            { name: "North Carolina" },
            { name: "North Dakota" },
            { name: "Ohio" },
            { name: "Oklahoma" },
            { name: "Oregon" },
            { name: "Pennsylvania" },
            { name: "Rhode Island" },
            { name: "South Carolina" },
            { name: "South Dakota" },
            { name: "Tennessee" },
            { name: "Texas" },
            { name: "Utah" },
            { name: "Vermont" },
            { name: "Virginia" },
            { name: "Washington" },
            { name: "West Virginia" },
            { name: "Wisconsin" },
            { name: "Wyoming" },
            { name: "District of Columbia" },
            { name: "Puerto Rico" },
            { name: "Guam" },
            { name: "American Samoa" },
            { name: "U.S. Virgin Islands" },
            { name: "Northern Mariana Islands" }
        ];
        await collection.drop().catch(err => console.log('Collection does not exist, skipping drop.'));
        
        await collection.insertMany(states);
    } finally {
        await client.close();
    }
}

main().catch(console.error);

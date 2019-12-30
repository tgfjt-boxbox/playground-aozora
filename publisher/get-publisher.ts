import admin from "firebase-admin";

admin.initializeApp();

const db = admin.firestore();

type Publisher = {
  name: string
}

async function main() {
  try {
    const ret = [] as Publisher[];
    const s = await db.collection("aozora_publishers").get();
    s.forEach((doc) => {
      ret.push(doc.data() as Publisher);
    });
    console.log(ret);
  } catch (err) {
    console.warn(err);
  }
}
main();

// to run:  $ts-node publisher/get-publisher.ts

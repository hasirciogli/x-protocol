import { newXProtocolClient } from "./packages/ts/src/packages/client";

(async () => {
  const client = newXProtocolClient("localhost", 8080);
  const response = await client.call({
    name: "hello",
    // token: "123456",
    payload: {
      username: "jhon",
      password: "doe1234",
    },
  });

  console.log(response);
})();

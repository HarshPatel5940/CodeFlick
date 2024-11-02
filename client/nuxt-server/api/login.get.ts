export default defineEventHandler(async (event) => {
  const BE = process.env.BE_URL;

  return {
    status: 307,
    data: {
      redirectURI: BE + "/api/auth/google/login?r=client",
    },
  };
});

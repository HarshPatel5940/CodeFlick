export default defineEventHandler(async (event) => {
  const BE = process.env.BE_URL

  const session = getCookie(event, 'session')

  if (!session) {
    return {
      status: 401,
      data: {
        message: 'Unauthorized! No Cookie Found.',
      },
    }
  }

  const response = await fetch(`${BE}/api/auth/session`, {
    headers: {
      Cookie: `session=${session}`,
    },
  })

  const body = await response.json()

  return {
    status: 200,
    data: body,
  }
})

export default defineEventHandler(async (event) => {
    const id = getRouterParam(event, 'id')
    if (!id) {
        throw createError({ statusCode: 400, statusMessage: 'Hotel ID is required' })
    }

    const backendUrl = process.env.BACKEND_API_URL || 'http://localhost:8080'
    try {
        return await $fetch(`${backendUrl}/api/hotels/${id}`)
    } catch (error: any) {
        throw createError({
            statusCode: error?.response?.status || 502,
            statusMessage: 'Unable to load hotel data from API'
        })
    }
})

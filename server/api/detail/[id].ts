import hotelDetailsData from '~/app/data/hotels/hotels_details.json'

export default defineEventHandler((event) => {
    const id = getRouterParam(event, 'id')
    const hotel = (hotelDetailsData as any[]).find(h => String(h.id) === String(id))

    if (!hotel) {
        return null
    }

    return hotel
})

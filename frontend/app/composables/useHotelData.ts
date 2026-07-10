import { joinURL } from 'ufo'

export const useHotelData = () => {
    const config = useRuntimeConfig()
    const baseURL = config.app.baseURL
    const defaultImage = joinURL(baseURL, 'data/images/_default.jpg')

    const handleImageError = (event: Event) => {
        const image = event.target as HTMLImageElement
        image.src = defaultImage
    }

    return {
        defaultImage,
        handleImageError
    }
}

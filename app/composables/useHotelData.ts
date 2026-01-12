import area1 from '~/data/areas/area_1.json'
import area2 from '~/data/areas/area_2.json'
import area3 from '~/data/areas/area_3.json'
import area4 from '~/data/areas/area_4.json'
import area5 from '~/data/areas/area_5.json'
import area6 from '~/data/areas/area_6.json'
import area7 from '~/data/areas/area_7.json'
import area8 from '~/data/areas/area_8.json'
// area_9 missing
import area10 from '~/data/areas/area_10.json'
import area11 from '~/data/areas/area_11.json'
import area12 from '~/data/areas/area_12.json'
import area13 from '~/data/areas/area_13.json'
import area14 from '~/data/areas/area_14.json'
// area_15 missing
import area16 from '~/data/areas/area_16.json'
// area_17 missing
import area18 from '~/data/areas/area_18.json'
import area19 from '~/data/areas/area_19.json'
import area20 from '~/data/areas/area_20.json'
import area21 from '~/data/areas/area_21.json'
import area22 from '~/data/areas/area_22.json'
import area23 from '~/data/areas/area_23.json'
import area24 from '~/data/areas/area_24.json'

import { joinURL } from 'ufo'

export const useHotelData = () => {
    const config = useRuntimeConfig()
    const baseURL = config.app.baseURL
    const defaultImage = joinURL(baseURL, 'data/images/_default.jpg')

    // Correct Mapping based on scripts/extract_city_names.js
    const cities = [
        { id: 1, name: '基隆' }, { id: 2, name: '台北' }, { id: 3, name: '新北' },
        { id: 4, name: '桃園' }, { id: 5, name: '新竹' }, { id: 6, name: '宜蘭' },
        { id: 7, name: '苗栗' }, { id: 8, name: '台中' }, { id: 10, name: '彰化' },
        { id: 11, name: '南投' }, { id: 12, name: '雲林' }, { id: 13, name: '嘉義' },
        { id: 14, name: '台南' }, { id: 16, name: '高雄' }, { id: 18, name: '屏東' },
        { id: 19, name: '花蓮' }, { id: 20, name: '台東' }, { id: 21, name: '澎湖' },
        { id: 22, name: '金門' }, { id: 23, name: '馬祖' }, { id: 24, name: '其他' }
    ]

    const regionCities = [
        {
            name: '北部',
            cities: [
                { id: 1, name: '基隆' }, { id: 2, name: '台北' }, { id: 3, name: '新北' },
                { id: 4, name: '桃園' }, { id: 5, name: '新竹' }, { id: 6, name: '宜蘭' }
            ]
        },
        {
            name: '中部',
            cities: [
                { id: 7, name: '苗栗' }, { id: 8, name: '台中' }, { id: 10, name: '彰化' },
                { id: 11, name: '南投' }, { id: 12, name: '雲林' }
            ]
        },
        {
            name: '南部',
            cities: [
                { id: 13, name: '嘉義' }, { id: 14, name: '台南' }, { id: 16, name: '高雄' },
                { id: 18, name: '屏東' }
            ]
        },
        {
            name: '東部',
            cities: [
                { id: 19, name: '花蓮' }, { id: 20, name: '台東' }
            ]
        },
        {
            name: '離島',
            cities: [
                { id: 21, name: '澎湖' }, { id: 22, name: '金門' }, { id: 23, name: '馬祖' }
            ]
        }
    ]

    const getCityDataById = (id: string | number) => {
        const areaMap: Record<string, any[]> = {
            '1': area1,
            '2': area2,
            '3': area3,
            '4': area4,
            '5': area5,
            '6': area6,
            '7': area7,
            '8': area8,
            // '9': area9,
            '10': area10,
            '11': area11,
            '12': area12,
            '13': area13,
            '14': area14,
            // '15': area15,
            '16': area16,
            // '17': area17,
            '18': area18,
            '19': area19,
            '20': area20,
            '21': area21,
            '22': area22,
            '23': area23,
            '24': area24
        }
        const data = areaMap[String(id)]
        return Array.isArray(data) ? data : []
    }

    const processCityData = (id: number, name: string, data: any[]) => {
        const safeData = Array.isArray(data) ? data : []
        return {
            id,
            name,
            hotels: safeData.slice(0, 6).map(item => {
                const id = item.link.split('n=')[1] || ''
                return {
                    id,
                    name: item.name,
                    image: joinURL(baseURL, `data/images/${id}.jpg`),
                    price: item.price,
                    address: item.address,
                    link: item.link
                }
            })
        }
    }

    // Updated selected IDs
    const selectedCitiesData = [
        processCityData(2, '台北', area2),
        processCityData(3, '新北', area3),
        processCityData(4, '桃園', area4), // Was 5
        processCityData(8, '台中', area8),
        processCityData(14, '台南', area14), // Was 13
        processCityData(16, '高雄', area16)  // Was 14
    ]

    const handleImageError = (e: Event) => {
        const img = (e.target as HTMLImageElement)
        img.src = defaultImage
    }

    return {
        cities,
        selectedCitiesData,
        regionCities,
        defaultImage,
        getCityDataById,
        handleImageError
    }
}

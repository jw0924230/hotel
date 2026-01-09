import area2 from '~/data/areas/area_2.json'
import area3 from '~/data/areas/area_3.json'
import area5 from '~/data/areas/area_5.json'
import area8 from '~/data/areas/area_8.json'
import area13 from '~/data/areas/area_13.json'
import area14 from '~/data/areas/area_14.json'

const defaultImage = '/data/images/default.jpg'

export const useHotelData = () => {
    const cities = [
        { id: 1, name: '基隆' }, { id: 2, name: '台北' }, { id: 3, name: '新北' },
        { id: 4, name: '宜蘭' }, { id: 5, name: '桃園' }, { id: 6, name: '新竹' },
        { id: 7, name: '苗栗' }, { id: 8, name: '台中' }, { id: 9, name: '彰化' },
        { id: 10, name: '南投' }, { id: 11, name: '雲林' }, { id: 12, name: '嘉義' },
        { id: 13, name: '台南' }, { id: 14, name: '高雄' }, { id: 15, name: '屏東' },
        { id: 16, name: '花蓮' }, { id: 17, name: '台東' }, { id: 18, name: '金門' },
        { id: 19, name: '澎湖' }, { id: 20, name: '馬祖' }, { id: 21, name: '其他' }
    ]

    const regionCities = [
        {
            name: '北部',
            cities: [
                { id: 1, name: '基隆' }, { id: 2, name: '台北' }, { id: 3, name: '新北' },
                { id: 5, name: '桃園' }, { id: 6, name: '新竹' }, { id: 4, name: '宜蘭' }
            ]
        },
        {
            name: '中部',
            cities: [
                { id: 7, name: '苗栗' }, { id: 8, name: '台中' }, { id: 9, name: '彰化' },
                { id: 10, name: '南投' }, { id: 11, name: '雲林' }
            ]
        },
        {
            name: '南部',
            cities: [
                { id: 12, name: '嘉義' }, { id: 13, name: '台南' }, { id: 14, name: '高雄' },
                { id: 15, name: '屏東' }
            ]
        },
        {
            name: '東部',
            cities: [
                { id: 16, name: '花蓮' }, { id: 17, name: '台東' }
            ]
        },
        {
            name: '離島',
            cities: [
                { id: 18, name: '金門' }, { id: 19, name: '澎湖' }, { id: 20, name: '馬祖' }
            ]
        }
    ]

    const getCityDataById = (id: string | number) => {
        const areaMap: Record<string, any[]> = {
            '2': area2,
            '3': area3,
            '5': area5,
            '8': area8,
            '13': area13,
            '14': area14
        }
        return areaMap[String(id)] || []
    }

    const processCityData = (id: number, name: string, data: any[]) => {
        return {
            id,
            name,
            hotels: data.slice(0, 6).map(item => {
                const id = item.link.split('n=')[1] || ''
                return {
                    id,
                    name: item.name,
                    image: `/data/images/${id}.jpg`,
                    price: item.price,
                    address: item.address,
                    link: item.link
                }
            })
        }
    }

    const selectedCitiesData = [
        processCityData(2, '台北', area2),
        processCityData(3, '新北', area3),
        processCityData(5, '桃園', area5),
        processCityData(8, '台中', area8),
        processCityData(13, '台南', area13),
        processCityData(14, '高雄', area14)
    ]

    return {
        cities,
        selectedCitiesData,
        regionCities,
        defaultImage,
        getCityDataById
    }
}

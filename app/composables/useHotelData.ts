import area2 from '~/data/areas/area_2.json'
import area3 from '~/data/areas/area_3.json'
import area5 from '~/data/areas/area_5.json'
import area8 from '~/data/areas/area_8.json'
import area13 from '~/data/areas/area_13.json'
import area14 from '~/data/areas/area_14.json'
import hotelDetailsData from '~/data/hotels/hotels_details.json'
import defaultImage from '~/data/images/default.jpg'

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


    // Load all images from data/images
    const hotelImages = import.meta.glob('~/data/images/*.jpg', { eager: true, import: 'default' })

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

    const getHotelById = (id: string | number) => {
        return (hotelDetailsData as any[]).find(h => String(h.id) === String(id))
    }

    const processCityData = (id: number, name: string, data: any[]) => {
        return {
            id,
            name,
            hotels: data.slice(0, 6).map(item => {
                const id = item.link.split('n=')[1] || ''
                // Construct the key that matches import.meta.glob keys
                // Glob patterns are relative to this file or root? 
                // In Nuxt/Vite, '~/data/images/*.jpg' usually resolves to keys like '/data/images/1234.jpg' or similar depending on config.
                // However, let's try to match loosely or inspect the keys if possible. 
                // A safe bet with eager glob is that keys are absolute or relative to project root alias.
                // Actually '~/data/images/*.jpg' will produce keys starting with /data/images/... if defined that way?
                // No, usually it preserves the pattern structure.
                // Let's assume the key format is like '/_nuxt/data/images/...' or just use a helper.

                // Better approach:
                // keys will be like "/data/images/1001.jpg" if we use absolute alias, 
                // or "../../data/images/1001.jpg" if relative.
                // Let's use string manipulation to find the image.

                // Using specific pattern for glob to ensure keys are predictable
                // const hotelImages = import.meta.glob('/data/images/*.jpg', { eager: true, import: 'default' })

                // Let's check if the ID exists in the images
                // We will try to find a key that ends with `/${id}.jpg`

                const imagePath = Object.keys(hotelImages).find(key => key.endsWith(`/${id}.jpg`))
                const imageUrl = imagePath ? (hotelImages[imagePath] as string) : defaultImage

                return {
                    id,
                    name: item.name,
                    image: imageUrl,
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

    const discountHotels = [
        { id: 1283, name: '友徠汽車旅館', city: '台北', price: '2480元起', description: '時尚客房', image: 'https://placehold.co/200x150/2c3e50/ffffff?text=友徠' },
        { id: 3527, name: '水雲端概念旅館', city: '台中', price: '2880元起', description: '精品商務房', image: 'https://placehold.co/200x150/e67e22/ffffff?text=水雲端' },
        { id: 2257, name: '紫晶彩繪汽車旅館', city: '新竹', price: '1780元起', description: '商務兩人房', image: 'https://placehold.co/200x150/27ae60/ffffff?text=紫晶' },
        { id: 3741, name: '挪威森林台中漫活館', city: '台中', price: '1780元起', description: '漫活商務房', image: 'https://placehold.co/200x150/d35400/ffffff?text=挪威森林' },
        { id: 3657, name: '紫禁城主題會館', city: '台中', price: '2580元起', description: '精品客棧', image: 'https://placehold.co/200x150/8e44ad/ffffff?text=紫禁城' },
        { id: 1861, name: '沐蘭精品台北館', city: '台北', price: '3110元起', description: '蝶舞套房', image: 'https://placehold.co/200x150/2c3e50/ffffff?text=沐蘭' }
    ]

    const regionHotels = {
        north: [
            { city: '台北', name: '友徠汽車旅館', type: '時尚客房', price: '2480元起' },
            { city: '台北', name: '沐蘭精品台北館', type: '蝶舞套房20:00入', price: '3110元起' },
            { city: '新北', name: '挪威森林新店館', type: '精品商務房', price: '1780元起' },
            { city: '新北', name: '紫晶彩繪汽車旅館', type: '尊爵商務', price: '2480元起' },
            { city: '桃園', name: '哈密瓜中壢館', type: '經典商務房', price: '1280元起' }
        ],
        central: [
            { city: '台中', name: '挪威森林台中漫活館', type: '漫活商務房', price: '1780元起' },
            { city: '台中', name: '雲河概念會館', type: '商務房', price: '1980元起' },
            { city: '台中', name: '水舞行館', type: '商務雙人房', price: '2980元起' },
            { city: '台中', name: '木木行館 Motel Lin', type: '商務雙人房', price: '2980元起' },
            { city: '台中', name: '紫禁城主題會館', type: '精品客棧', price: '2580元起' }
        ],
        south: [
            { city: '雲林', name: '御品王朝旅店', type: '商務客房', price: '2100元起' },
            { city: '南投', name: '海悅時尚汽車旅館', type: '典雅商務', price: '2280元起' },
            { city: '台南', name: '歐薇汽車旅館', type: '時尚套房', price: '2280元起' },
            { city: '高雄', name: '花鄉戀館大昌店', type: '簡約精品', price: '1480元起' },
            { city: '高雄', name: '欣國賓商務旅館', type: '標準雙人房', price: '2200元起' }
        ]
    }

    // Detailed Data for Inquiry Page
    const hotelDetails = {
        7: {
            id: 7,
            name: '麗景經典汽車旅館',
            city: '彰化',
            area: '鹿港鎮',
            address: '彰化縣鹿港鎮長安路58號',
            phone: '04-777-6666',
            website: 'http://www.easehotel.com.tw',
            email: 'regionmotel@yahoo.com.tw',
            image: 'https://placehold.co/600x400/2c3e50/ffffff?text=麗景經典',
            desc: `麗景經典汽車旅館位於鹿港鎮長安路，鄰近鹿港老街。
顛覆傳統！只給您最好的享受-義式研磨咖啡
最新引進花灑SPA 等您鑑賞
免費欣賞影音劇院`,
            roomTypes: [
                { name: '經典商務房B', price: '1680元起', img: 'https://placehold.co/120x90/333/fff?text=商務B' },
                { name: '豪華型', price: '2480元起', img: 'https://placehold.co/120x90/555/fff?text=豪華' },
                { name: '標準雙人房', price: '2200元起', img: 'https://placehold.co/120x90/777/fff?text=雙人' }
            ],
            news: [
                '最新消息：顛覆傳統！只給您最好的享受-義式研磨咖啡',
                '最新消息：最新引進花灑SPA 等您鑑賞',
                '最新消息：免費欣賞影音劇院'
            ]
        }
    }

    // Expanded list for search results (simulating database)
    const allHotels = [
        { id: 1132, name: '彰化桂冠精品旅館', city: '彰化', price: '2480元起', image: 'https://placehold.co/240x160/e74c3c/ffffff?text=彰化桂冠', desc: '持QK卡優惠：休息折100元 ； 住宿折300元', type: '汽車旅館' },
        { id: 8, name: '凱登汽車旅館', city: '彰化', price: '住宿折200元', image: 'https://placehold.co/240x160/2c3e50/ffffff?text=凱登', desc: '持QK卡優惠：住宿折200元', type: '汽車旅館' },
        { id: 1331, name: '第八月台汽車旅館', city: '彰化', price: '休息折100元', image: 'https://placehold.co/240x160/8e44ad/ffffff?text=第八月台', desc: '仲夏自由FUN 暑期優惠活動', type: '汽車旅館' },
        { id: 3679, name: '歐遊連鎖精品旅館-彰化館', city: '彰化', price: '詳全文', image: 'https://placehold.co/240x160/3498db/ffffff?text=歐遊彰化', desc: '全台連鎖精品旅館', type: '汽車旅館' },
        { id: 831, name: '里昂晶鑽汽車旅館', city: '彰化', price: '詳全文', image: 'https://placehold.co/240x160/16a085/ffffff?text=里昂晶鑽', desc: '精緻休閒住宿', type: '汽車旅館' },
        { id: 28, name: '加州汽車旅館', city: '彰化', price: '詳全文', image: 'https://placehold.co/240x160/f39c12/ffffff?text=加州', desc: '舒適休憩空間', type: '汽車旅館' },
        { id: 14, name: '富晴汽車旅館', city: '彰化', price: '詳全文', image: 'https://placehold.co/240x160/c0392b/ffffff?text=富晴', desc: '平價優質選擇', type: '汽車旅館' },
        { id: 1283, name: '友徠汽車旅館', city: '台北', price: '2480元起', image: 'https://placehold.co/240x160/2c3e50/ffffff?text=友徠', desc: '時尚客房', type: '汽車旅館' },
        { id: 3741, name: '挪威森林台中漫活館', city: '台中', price: '1780元起', image: 'https://placehold.co/240x160/d35400/ffffff?text=挪威森林', desc: '漫活商務房', type: '汽車旅館' }
    ]

    const news = [
        { id: 1, title: '普發1萬元 快閃優惠票卷 限時大搶購', city: '台北', store: '雅柏精緻汽車旅館', tag: '最新' },
        { id: 2, title: '凱米變中颱發布海警', city: '新北', store: '欣埔山莊', tag: '快訊' },
        { id: 3, title: '114年8/29 七夕限時活動', city: '桃園', store: '人客旅館', tag: '活動' },
        { id: 4, title: '親子互動商品，累積點數免費兌換哦~', city: '新北', store: '133汽車旅館', tag: '優惠' },
        { id: 5, title: '周六住宿不加價優惠中', city: '雲林', store: '采弘汽車商務旅館', tag: '優惠' }
    ]

    return {
        cities,
        discountHotels,
        regionHotels,
        allHotels,
        hotelDetails,
        news,
        selectedCitiesData,
        regionCities,
        defaultImage,
        getCityDataById,
        getHotelById
    }
}

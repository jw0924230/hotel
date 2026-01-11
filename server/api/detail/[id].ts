import fs from 'fs'
import path from 'path'

export default defineEventHandler((event) => {
    const id = getRouterParam(event, 'id')
    if (!id) return null

    const filePath = path.resolve(process.cwd(), 'data_json/hotels/details', `${id}.json`)

    if (!fs.existsSync(filePath)) {
        return null
    }

    try {
        const data = fs.readFileSync(filePath, 'utf-8')
        return JSON.parse(data)
    } catch (e) {
        console.error(`Error reading hotel data for ID ${id}:`, e)
        return null
    }
})

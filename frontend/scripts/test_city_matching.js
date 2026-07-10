
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const detailsDir = path.resolve(__dirname, '../data_json/hotels/details');

const cities = [
    { id: 1, name: '基隆' }, { id: 2, name: '台北' }, { id: 3, name: '新北' },
    { id: 4, name: '宜蘭' }, { id: 5, name: '桃園' }, { id: 6, name: '新竹' },
    { id: 7, name: '苗栗' }, { id: 8, name: '台中' }, { id: 9, name: '彰化' },
    { id: 10, name: '南投' }, { id: 11, name: '雲林' }, { id: 12, name: '嘉義' },
    { id: 13, name: '台南' }, { id: 14, name: '高雄' }, { id: 15, name: '屏東' },
    { id: 16, name: '花蓮' }, { id: 17, name: '台東' }, { id: 18, name: '金門' },
    { id: 19, name: '澎湖' }, { id: 20, name: '馬祖' }, { id: 21, name: '其他' }
];

const files = fs.readdirSync(detailsDir).filter(f => f.endsWith('.json'));

let matched = 0;
let failed = 0;
const failures = [];

files.forEach(file => {
    const content = fs.readFileSync(path.join(detailsDir, file), 'utf-8');
    const hotel = JSON.parse(content);

    if (!hotel.address) {
        failures.push({ id: hotel.id, address: 'N/A', reason: 'No address' });
        failed++;
        return;
    }

    const addr = hotel.address;
    // Mimic the logic in the component
    // const name = addr.substring(0, 3).replace('台', '臺'); 

    const found = cities.find(c => addr.includes(c.name));

    if (found) {
        matched++;
    } else {
        failures.push({ id: hotel.id, address: addr, reason: 'No match found' });
        failed++;
    }
});

console.log(`Total Files: ${files.length}`);
console.log(`Matched: ${matched}`);
console.log(`Failed: ${failed}`);

if (failed > 0) {
    console.log('\nTop 20 Failures:');
    failures.slice(0, 20).forEach(f => console.log(`${f.id}: ${f.address}`));
}

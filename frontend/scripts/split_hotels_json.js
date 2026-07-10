
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const sourceFile = path.resolve(__dirname, '../app/data/hotels/hotels_details.json');
const outputDir = path.resolve(__dirname, '../app/data/hotels/details');

if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
}

console.log('Reading source file...');
const data = JSON.parse(fs.readFileSync(sourceFile, 'utf-8'));

console.log(`Found ${data.length} hotels. Splitting...`);

let count = 0;
for (const hotel of data) {
    if (hotel.id) {
        fs.writeFileSync(path.join(outputDir, `${hotel.id}.json`), JSON.stringify(hotel, null, 2));
        count++;
    }
}

console.log(`Successfully split ${count} files to ${outputDir}`);

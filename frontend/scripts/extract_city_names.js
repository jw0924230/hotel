
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const areasDir = path.resolve(__dirname, '../app/data/areas');
const files = fs.readdirSync(areasDir).filter(f => f.startsWith('area_') && f.endsWith('.json'));

const mapping = [];

files.forEach(file => {
    // Extract ID from filename area_X.json
    const id = parseInt(file.replace('area_', '').replace('.json', ''));

    try {
        const content = fs.readFileSync(path.join(areasDir, file), 'utf-8');
        const data = JSON.parse(content);

        if (Array.isArray(data) && data.length > 0) {
            // Assume the first item has the correct "area" (city) name
            // The JSON structure from previous view was: { "area": "马祖", ... }
            const name = data[0].area;
            mapping.push({ id, name });
        } else {
            mapping.push({ id, name: 'Unknown (Empty)' });
        }
    } catch (e) {
        mapping.push({ id, name: 'Error reading file' });
    }
});

// Sort by ID
mapping.sort((a, b) => a.id - b.id);

console.log('Detected Mapping:');
mapping.forEach(m => {
    console.log(`{ id: ${m.id}, name: '${m.name}' },`);
});

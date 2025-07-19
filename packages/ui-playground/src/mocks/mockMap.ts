import { createMockBlock } from './mockBlocks';

// List of possible block types
const blockTypes = ['blank', 'castle', 'king', 'mountain', 'soldier'];

// Function to create a 20x20 map with random block types
export function generateMockMap(): ReturnType<typeof createMockBlock>[][] {
    const width = 20;
    const height = 20;
    const map = [];

    for (let y = 0; y < height; y++) {
        const row = [];
        for (let x = 0; x < width; x++) {
            // Randomly select a block type
            const blockType = blockTypes[Math.floor(Math.random() * blockTypes.length)];
            const owner = blockType === 'blank' || blockType === 'mountain' ? 0 : (Math.random() > 0.5 ? 1 : 2);
            // Create and add the block to the row
            row.push(createMockBlock(blockType, y * width + x, owner));
        }
        map.push(row);
    }

    return map;
}


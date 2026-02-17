import { useInfiniteQuery } from '@tanstack/react-query';

export interface BrainrotItem {
    id: string;
    type: 'image' | 'video';
    url: string;
    title?: string;
    aspectRatio: number;
}

const MOCK_IMAGES = [
    'https://pbs.twimg.com/media/GFPI8LKXQAAis33.jpg',
    'https://ih1.redbubble.net/image.5681776202.5017/st,small,507x507-pad,600x600,f8f8f8.jpg',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRj8_eBhg2oB0B05Ke2OptVopeUnV5nPCEh9Q&s',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSvkrJFE1in9bcHG0IAsLVqP_AVYIbUlmtbGg&s',
    'https://media.tenor.com/kIjcbnIniMsAAAAe/sukuna.png',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTIQ40UACBhEOuo_xTJAAq3OCP08GVt7n3Dwg&s',
    'https://media.tenor.com/UzBRE6feDAEAAAAe/he-made-a-statement-so-trash-even-his-gang-clowned-him-dog.png',
    'https://media.tenor.com/D8NPCEKVngsAAAAe/joker-why-so-serious.png',
    'https://i.pinimg.com/736x/bf/6e/29/bf6e296386c67b027cd3d234e3c6efa4.jpg',
    'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ6Gtano8OnHkaq-6hB7G-ts9KVw99Ye5ReSw&s',
];

async function fetchBrainrot(page: number): Promise<{ items: BrainrotItem[]; nextPage: number }> {
    // No simulated delay for better UX

    const items: BrainrotItem[] = Array.from({ length: 10 }).map((_, i) => {
        const index = (page * 10 + i) % MOCK_IMAGES.length;
        return {
            id: `brainrot-${page}-${i}-${Date.now()}`,
            type: 'image',
            url: MOCK_IMAGES[index],
            aspectRatio: 1,
        };
    });

    return {
        items,
        nextPage: page + 1,
    };
}

export function useBrainrotFeed() {
    return useInfiniteQuery({
        queryKey: ['brainrot'],
        queryFn: ({ pageParam }) => fetchBrainrot(pageParam),
        initialPageParam: 0,
        getNextPageParam: (lastPage) => lastPage.nextPage,
    });
}

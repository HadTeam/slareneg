import { createSignal, onMount, Show } from 'solid-js';
import { useSearchParams } from '@solidjs/router';
import Board from '../components/Board';
import TopBar from '../components/TopBar';
import { UploadPanel } from '../components/UploadPanel';
import type { ExportedMap } from '@slareneg/shared-types';

function TestBoard() {
  const [search, setSearch] = useSearchParams();
  const [mapData, setMapData] = createSignal<ExportedMap | null>(null);
  const [loading, setLoading] = createSignal(false);

  // Function to fetch random map
  const fetchRandomMap = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/map/random');
      if (response.ok) {
        const data = await response.json();
        setMapData(data);
        // Update URL without reloading
        setSearch({ random: '1' });
      } else {
        console.error('Failed to fetch random map:', response.status);
      }
    } catch (error) {
      console.error('Failed to fetch random map:', error);
    } finally {
      setLoading(false);
    }
  };

  // Fetch random map if query param is set on mount
  onMount(async () => {
    if (search.random === '1') {
      await fetchRandomMap();
    }
  });

  const handleGenerateRandom = async () => {
    await fetchRandomMap();
  };

  return (
    <div class="h-screen flex flex-col m-0 p-0 overflow-hidden">
      <TopBar />
      <Show
        when={mapData()}
        fallback={
          <div class="flex-1 flex flex-col items-center justify-center p-8">
            <div class="w-full max-w-xl mb-6">
              <UploadPanel onLoad={setMapData} />
            </div>
            <button
              onClick={handleGenerateRandom}
              disabled={loading()}
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading() ? 'Loading...' : 'Generate Random Map'}
            </button>
          </div>
        }
      >
        <div class="flex-1 flex flex-col">
          <Board blocks={mapData()!.blocks} size={mapData()!.size} />
          <div class="p-4 flex justify-center">
            <button
              onClick={() => {
                setMapData(null);
                setSearch({});
              }}
              class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors"
            >
              Clear Map
            </button>
          </div>
        </div>
      </Show>
    </div>
  );
}

export default TestBoard;

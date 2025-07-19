import { createSignal, onMount, Show } from 'solid-js';
import { useSearchParams, useNavigate } from '@solidjs/router';
import Board from '../components/Board';
import TopBar from '../components/TopBar';
import { UploadPanel } from '../components/UploadPanel';
import type { ExportedMap } from '@slareneg/shared-types';

function TestBoard() {
  const [search] = useSearchParams();
  const navigate = useNavigate();
  const [mapData, setMapData] = createSignal<ExportedMap | null>(null);

  // Fetch random map if query param is set
  onMount(async () => {
    if (search.random === '1') {
      try {
        const response = await fetch('/api/map/random');
        if (response.ok) {
          const data = await response.json();
          setMapData(data);
        }
      } catch (error) {
        console.error('Failed to fetch random map:', error);
      }
    }
  });

  const handleGenerateRandom = () => {
    navigate('?random=1');
    window.location.reload(); // Force reload to trigger the fetch
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
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Generate Random Map
            </button>
          </div>
        }
      >
        <Board blocks={mapData()!.blocks} size={mapData()!.size} />
      </Show>
    </div>
  );
}

export default TestBoard;

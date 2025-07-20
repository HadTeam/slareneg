import { createSignal, onMount, Show } from 'solid-js';
import { useSearchParams } from '@solidjs/router';
import Board from '../components/Board';
import Sidebar from '../components/Sidebar';
import { UploadPanel } from '../components/UploadPanel';
import type { ExportedMap } from '@slareneg/shared-types';

function TestBoard() {
  const [search, setSearch] = useSearchParams();
  const [mapData, setMapData] = createSignal<ExportedMap | null>(null);
  const [loading, setLoading] = createSignal(false);
  const [boardRef, setBoardRef] = createSignal<{ fitToView: () => void } | null>(null);

  // Function to fetch random map
  const fetchRandomMap = async () => {
    setLoading(true);
    // Clear existing map data to ensure UI updates
    setMapData(null);
    try {
      // Add timestamp to force a new request and bypass any caching
      const response = await fetch(`/api/map/random?t=${Date.now()}`);
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

  const clearMap = () => {
    setMapData(null);
    setSearch({});
  };

  return (
    <div class="h-screen flex overflow-hidden">
      <Sidebar
        onClearMap={clearMap}
        onFitToView={() => boardRef()?.fitToView()}
        onGenerateRandom={handleGenerateRandom}
        loading={loading()}
        hasMap={!!mapData()}
      />
      <Show
        when={mapData()}
        fallback={
          <div class="flex-1 flex flex-col items-center justify-center p-8">
            <div class="w-full max-w-xl">
              <UploadPanel onLoad={setMapData} />
              <div class="mt-4 text-center">
                <button
                  onClick={handleGenerateRandom}
                  disabled={loading()}
                  class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading() ? 'Loading...' : 'Generate Random Map'}
                </button>
              </div>
            </div>
          </div>
        }
      >
        <div class="flex-1 flex flex-col">
          <Board 
            blocks={mapData()!.blocks} 
            size={mapData()!.size} 
            onBoardRef={setBoardRef}
          />
        </div>
      </Show>
    </div>
  );
}

export default TestBoard;

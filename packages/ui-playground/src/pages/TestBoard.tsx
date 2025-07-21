import { createSignal, onMount, Show } from 'solid-js';
import { useSearchParams } from '@solidjs/router';
import Board from '../components/Board';
import Sidebar from '../components/Sidebar';
import { UploadPanel } from '../components/UploadPanel';
import type { ExportedMap } from '@slareneg/shared-types';
import { createTestBoard } from '../mocks/mockBlocks';

function TestBoard() {
  const [search, setSearch] = useSearchParams();
  const [mapData, setMapData] = createSignal<ExportedMap | null>(null);
  const [loading, setLoading] = createSignal(false);
  const [boardRef, setBoardRef] = createSignal<{ fitToView: () => void } | null>(null);


  const fetchRandomMap = async () => {
    setLoading(true);
    setMapData(null);
    try {
      const response = await fetch(`/api/map/random?t=${Date.now()}`);
      if (response.ok) {
        const data = await response.json();
        setMapData(data);
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

  onMount(async () => {
    if (search.random === '1') {
      await fetchRandomMap();
    } else if (search.mock === '1') {
      useMockData();
    }
  });

  const handleGenerateRandom = async () => {
    await fetchRandomMap();
  };

  const clearMap = () => {
    setMapData(null);
    setSearch({});
  };

  const useMockData = () => {
    setMapData(null);
    setTimeout(() => {
      const mockBlocks = createTestBoard(20, 20);
      const plainBlocks = mockBlocks.map(row => 
        row.map(block => ({
          meta: block.meta(),
          owner: block.owner(),
          num: block.num()
        }))
      );
      const mockMap: ExportedMap = {
        size: { width: 20, height: 20 },
        info: {
          id: `mock-map-${Date.now()}`,
          name: 'Mock Test Map',
          desc: `A test map with mock data (generated at ${new Date().toLocaleTimeString()})`
        },
        blocks: plainBlocks as any
      };
      setMapData(mockMap);
      setSearch({ mock: '1' });
    }, 0);
  };

  return (
    <div class="h-screen flex overflow-hidden">
      <Sidebar
        onClearMap={clearMap}
        onFitToView={() => boardRef()?.fitToView()}
        onGenerateRandom={handleGenerateRandom}
        onUseMockData={useMockData}
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

import { A } from '@solidjs/router';

interface SidebarProps {
  onClearMap?: () => void;
  onFitToView?: () => void;
  onGenerateRandom?: () => void;
  loading?: boolean;
  hasMap?: boolean;
}

function Sidebar(props: SidebarProps) {
  const { onClearMap, onFitToView, onGenerateRandom, loading = false, hasMap = false } = props;

  return (
    <div class="w-64 bg-gray-100 border-r border-gray-300 flex flex-col h-full">
      <div class="p-4 border-b border-gray-300">
        <h2 class="text-lg font-semibold mb-1">Test Board - Slareneg</h2>
        <p class="text-xs text-gray-600">Click to select • Drag to move • Scroll to zoom</p>
      </div>
      
      <div class="flex-1 p-4 space-y-4">
        <A 
          href="/" 
          class="block px-3 py-2 bg-gray-600 text-white no-underline rounded text-sm hover:bg-gray-700 text-center transition-colors"
        >
          ← Back to Home
        </A>
        
        {hasMap && (
          <>
            <button
              onClick={onFitToView}
              class="w-full px-3 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 transition-colors"
            >
              Fit to View
            </button>
            
            <button
              onClick={onClearMap}
              class="w-full px-3 py-2 bg-gray-600 text-white rounded text-sm hover:bg-gray-700 transition-colors"
            >
              Clear Map
            </button>
          </>
        )}
        
        {!hasMap && onGenerateRandom && (
          <button
            onClick={onGenerateRandom}
            disabled={loading}
            class="w-full px-3 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Loading...' : 'Generate Random Map'}
          </button>
        )}
      </div>
      
      <div class="p-4 border-t border-gray-300 text-xs text-gray-600">
        <h3 class="font-semibold mb-2">Legend:</h3>
        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <span class="text-xl">🏛</span> Castle
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xl">♔</span> King
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xl">⛰</span> Mountain
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xl">◆</span> Soldier
          </div>
        </div>
      </div>
    </div>
  );
}

export default Sidebar;

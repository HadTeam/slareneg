import { A } from '@solidjs/router';

function TopBar() {
  return (
    <div class="py-2 bg-gray-100 border-b border-gray-300 flex items-center gap-3 px-3">
      <A href="/" class="px-3 py-1 bg-gray-600 text-white no-underline rounded text-sm hover:bg-gray-700">
        ← Back to Home
      </A>
      <div>
        <h2 class="m-0 text-sm font-semibold">Test Board - Slareneg</h2>
        <p class="m-0 text-xs text-gray-600">Click to select • Drag to move • Scroll to zoom</p>
      </div>
    </div>
  );
}

export default TopBar;

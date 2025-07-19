import { createSignal } from 'solid-js';
import type { Component } from 'solid-js';
import type { ExportedMap } from '@slareneg/shared-types';

interface UploadPanelProps {
  onLoad: (map: ExportedMap) => void;
}

export const UploadPanel: Component<UploadPanelProps> = (props) => {
  const [isDragging, setIsDragging] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  const handleFile = async (file: File) => {
    setError(null);
    
    try {
      const text = await file.text();
      const map: ExportedMap = JSON.parse(text);
      
      // Validate basic fields
      if (!map.size || typeof map.size.width !== 'number' || typeof map.size.height !== 'number') {
        throw new Error('Invalid map: missing or invalid size field');
      }
      
      if (!Array.isArray(map.blocks)) {
        throw new Error('Invalid map: blocks must be an array');
      }
      
      if (map.blocks.length === 0) {
        throw new Error('Invalid map: blocks array cannot be empty');
      }
      
      // Call parent callback
      props.onLoad(map);
    } catch (err) {
      if (err instanceof SyntaxError) {
        setError('Invalid JSON file');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Unknown error occurred');
      }
    }
  };

  const handleFileInput = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
      handleFile(file);
    }
  };

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const file = e.dataTransfer?.files[0];
    if (file && (file.type === 'application/json' || file.name.endsWith('.json'))) {
      handleFile(file);
    } else {
      setError('Please drop a JSON file');
    }
  };

  return (
    <div
      class={`border-2 border-dashed rounded-lg p-6 text-gray-400 text-center transition-colors ${
        isDragging() ? 'border-blue-400 bg-blue-50' : 'border-gray-300'
      }`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <div class="mb-4">
        <svg
          class="mx-auto h-12 w-12 text-gray-400"
          stroke="currentColor"
          fill="none"
          viewBox="0 0 48 48"
          aria-hidden="true"
        >
          <path
            d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
      
      <label for="file-upload" class="cursor-pointer">
        <span class="text-sm font-medium text-blue-600 hover:text-blue-500">
          Upload a JSON map file
        </span>
        <input
          id="file-upload"
          type="file"
          accept="application/json"
          class="sr-only"
          onChange={handleFileInput}
        />
      </label>
      
      <p class="mt-2 text-sm">or drag and drop</p>
      
      {error() && (
        <p class="mt-2 text-sm text-red-600">{error()}</p>
      )}
    </div>
  );
};

import { useState } from "react";

export default function () {
  const [count, setCount] = useState(0);

  return (
    <div className="flex flex-col justify-center items-center h-[90vh]">
      <h1>Vite + React</h1>
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        onClick={() => setCount((count) => count + 1)}
      >
        count is {count}
      </button>
    </div>
  );
}

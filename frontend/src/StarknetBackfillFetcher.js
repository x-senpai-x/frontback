import React, { useState } from 'react';

const StarknetBackfillFetcher = () => {
  const [fromBlock, setFromBlock] = useState('');
  const [toBlock, setToBlock] = useState('');
  const [rpcUrl, setRpcUrl] = useState('');
  const [outputFile, setOutputFile] = useState('block_details.csv');
  const [transactionHashFlag, setTransactionHashFlag] = useState(false);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);

  const handleBackfill = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/backfill', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          fromBlock: Number(fromBlock),
          toBlock: Number(toBlock),
          rpcUrl: rpcUrl,
          outputFile: outputFile,
          transactionHashFlag: transactionHashFlag,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error ${response.status}`);
      }

      setSuccessMessage('Backfill successful!');
      setError(null);
    } catch (err) {
      setError(err.message);
      setSuccessMessage(null);
    }
  };

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Backfill StarkNet Block Data</h2>

      <div className="space-y-4">
        <div>
          <label htmlFor="fromBlock" className="block font-medium mb-1">
            From Block Number
          </label>
          <input
            type="number"
            id="fromBlock"
            className="w-full px-3 py-2 border rounded-md"
            value={fromBlock}
            onChange={(e) => setFromBlock(e.target.value)}
            placeholder="Enter starting block number"
          />
        </div>

        <div>
          <label htmlFor="toBlock" className="block font-medium mb-1">
            To Block Number
          </label>
          <input
            type="number"
            id="toBlock"
            className="w-full px-3 py-2 border rounded-md"
            value={toBlock}
            onChange={(e) => setToBlock(e.target.value)}
            placeholder="Enter ending block number"
          />
        </div>

        <div>
          <label htmlFor="rpcUrl" className="block font-medium mb-1">
            RPC URL
          </label>
          <input
            type="text"
            id="rpcUrl"
            className="w-full px-3 py-2 border rounded-md"
            value={rpcUrl}
            onChange={(e) => setRpcUrl(e.target.value)}
            placeholder="Enter the RPC URL"
          />
        </div>

        <div>
          <label htmlFor="outputFile" className="block font-medium mb-1">
            Output File
          </label>
          <input
            type="text"
            id="outputFile"
            className="w-full px-3 py-2 border rounded-md"
            value={outputFile}
            onChange={(e) => setOutputFile(e.target.value)}
            placeholder="Enter the output file name"
          />
        </div>

        <div className="flex items-center">
          <input
            type="checkbox"
            id="transactionHashFlag"
            checked={transactionHashFlag}
            onChange={(e) => setTransactionHashFlag(e.target.checked)}
          />
          <label htmlFor="transactionHashFlag" className="ml-2">
            Include Transaction Hashes
          </label>
        </div>

        <button
          className="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md"
          onClick={handleBackfill}
        >
          Fetch Block Data
        </button>
      </div>

      {successMessage && (
        <div className="mt-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative">
          <strong>{successMessage}</strong>
        </div>
      )}

      {error && (
        <div className="mt-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative">
          <strong>Error:</strong> {error}
        </div>
      )}
    </div>
  );
};

export default StarknetBackfillFetcher;
import React, { MouseEventHandler, useState } from 'react';

type CopyButtonProps = {
  textToCopy: string;
};

export const CopyButton: React.FC<CopyButtonProps> = ({ textToCopy }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy:MouseEventHandler<HTMLButtonElement> = async (e) => {
    e.preventDefault()
    try {
      await navigator.clipboard.writeText(textToCopy);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000); // Volta ao estado normal depois de 2s
    } catch (err) {
      console.error('Erro ao copiar!', err);
    }
  };

  return (
    <button onClick={handleCopy} className="px-4 py-2 bg-blue-500 text-white rounded">
      {copied ? 'Copiado!' : 'Copiar'}
    </button>
  );
};

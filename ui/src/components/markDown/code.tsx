import React from 'react';

const Code = () => {
  return (
    <div style={{ height: 420 }}>
      <MonacoEditor
        height='100%'
        language={getBaseLanguageId(data?.program_language || 'plaintext')}
        value={editorValue}
        theme='vs-dark'
        options={{
          readOnly: true,
          minimap: { enabled: false },
          fontSize: 14,
          scrollBeyondLastLine: false,
          wordWrap: 'on',
          lineNumbers: 'on',
          glyphMargin: false,
          folding: false,
          overviewRulerLanes: 0,
          guides: {
            indentation: true,
            highlightActiveIndentation: true,
            highlightActiveBracketPair: false,
          },
          renderLineHighlight: 'none',
          cursorStyle: 'line',
          cursorBlinking: 'solid',
          cursorWidth: 0,
          contextmenu: false,
          selectionHighlight: false,
          selectOnLineNumbers: false,
          occurrencesHighlight: 'off',
          links: false,
          hover: { enabled: false },
          codeLens: false,
          dragAndDrop: false,
          mouseWheelZoom: false,
          accessibilitySupport: 'off',
          bracketPairColorization: { enabled: false },
          matchBrackets: 'never',
        }}
        onMount={(editor) => {
          editorRef.current = editor;
          setEditorReady(true);
          // 隐藏光标
          const editorDom = editor.getDomNode();
          if (editorDom) {
            const style = document.createElement('style');
            style.innerHTML = `.monaco-editor .cursor { display: none !important; }`;
            editorDom.appendChild(style);
          }
        }}
      />
    </div>
  );
};

export default Code;

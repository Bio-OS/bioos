/**
 *
 * Copyright 2023 Beijing Volcano Engine Technology Ltd.
 * Copyright 2023 Guangzhou Laboratory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Icon from 'components/Icon';

const icons = {
  go: <Icon glyph="editor-go" className="arco-icon basic-editor-file-icon" />,
  java: (
    <Icon glyph="editor-java" className="arco-icon basic-editor-file-icon" />
  ),
  js: (
    <Icon glyph="editor-nodejs" className="arco-icon basic-editor-file-icon" />
  ),
  py: (
    <Icon glyph="editor-python" className="arco-icon basic-editor-file-icon" />
  ),
  ts: (
    <Icon glyph="editor-nodejs" className="arco-icon basic-editor-file-icon" />
  ),
  wdl: <Icon glyph="wdl" className="arco-icon basic-editor-file-icon" />,
};

export default function EditorFileIcon({
  suffix,
  dir,
}: {
  suffix?: string;
  dir?: boolean;
}) {
  if (dir) {
    return (
      <Icon
        glyph="editor-folder"
        className="arco-icon basic-editor-file-icon"
      />
    );
  }

  return (
    icons[suffix] || (
      <Icon glyph="editor-file" className="arco-icon basic-editor-file-icon" />
    )
  );
}

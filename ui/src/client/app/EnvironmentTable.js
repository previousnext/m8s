import React from 'react';

import Options from './Options';
import UIs from './UIs';

export default function EnvironmentTable (props) {
  return (
    <table className="js-table--responsive">
      <thead>
        <tr>
          <th>Domain</th>
          <th>Logs</th>
          <th>Console</th>
          <th>UIs</th>
        </tr>
      </thead>
      <tbody>
        {props.envs.filter((env) => {
          return props.search === ''
            ? true
            : env.name.includes(props.search)
        })
        .map((env) => {
          return (
            <tr key={env.name}>
              <td className="table--enlarged">
                <a href={`//${env.domain}`}>{env.name}</a>
              </td>
              <td>
                <Options operation="logs" name={env.name} containers={env.containers} />
              </td>
              <td>
                <Options operation="shell" name={env.name} containers={env.containers} />
              </td>
              <td>
                { /* @todo, Consolidate with the Options component. */ }
                <UIs name={env.name} base_url={`//${env.domain}`} />
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}

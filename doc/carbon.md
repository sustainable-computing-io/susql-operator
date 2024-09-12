# Carbon Dioxide Emission Estimation

There are two currently CO2 emission calculation methods available, and an additional method is under development.

"Out-of-the-box" SusQL reports an estimated CO2 emission value for all measured workloads using the `static` method described below.

### `static` Method
- The default CO2 emission calculation method is used when the `CARBON-METHOD` ConfigMap value is set to `static`
which simply uses a grams of CO2 per Joule 
of electricity consumed coefficient to calculate grams of CO2 emitted. This value is user tunable by
modifying the `CARBON-INTENSITY` ConfigMap value.  The default value is based on
[US EPA](https://www.epa.gov/energy/greenhouse-gases-equivalencies-calculator-calculations-and-references)
data.

### `simpledynamic` Method
- The `simpledynamic` method periodically queries the carbon intensity value for a user specified location to provide a more accurate estimation of CO2 emission.
The ConfigMap user configurable items are:
  - `CARBON-METHOD` - The `simpledynamic` method is enabled when this set to `simpledynamic`.
  - `CARBON-INTENSITY` - This value is set automatically with `simpledynamic`. User specified values will be overwritten.
  - `CARBON-INTENSITY-URL` - Specifies a web API that returns carbon intensity. The default value works as of the date of this writing.
  - `CARBON-LOCATION` - The location identifiers are defined by the API provider, but can be selected by he user.
  - `CARBON-QUERY-RATE` - Interval in seconds at which the carbon intensity is queried. The data available from the source is updated less than hourly, so an interval of greater than 3600 seconds is recommended.
  - `CARBON-QUERY-FILTER` - When the return value is embedded in a JSON object, this specification enables the extraction of the data. The default value matches the default provider.
  - `CARBON-QUERY-CONV-2J` - If the carbon data provider does not provide data in the standard "grams of CO2 per Joule" then this factor can be specified to normalize the units displayed. The default value ensures that the default provider data is in the correct unit.

### `sdk` Method
- The third method `sdk` is still under development, and will offer an integration with the Green Software Foundation's [Carbon Aware SDK](https://github.com/Green-Software-Foundation/carbon-aware-sdk).

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


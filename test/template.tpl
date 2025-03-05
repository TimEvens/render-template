{{- /*
   ---------------------------------------------------------------------------
   * Define variables that are used in this template
   *    NOTE: Only define variables that need to be generated from the
   *          values file(s).  For example, concatenate values into
   *          a single variable that is then reused.
   ---------------------------------------------------------------------------
*/ -}}

{{- $VRF_NAME := cat "DI-CUST-" .customer.id | nospace -}}
{{- $VASI_IPS := vasiLinkIps "100.64.0.1/22" 10 -}}

{{- /*
   ---------------------------------------------------------------------------
   * START of template
   ---------------------------------------------------------------------------
*/ -}}

Object Name {{ .object.name }}
Second name {{ .object.name2 }}

Customer Id {{ .customer.id }}

VRF Name {{ $VRF_NAME }}

Link IPs:
    Left {{ index $VASI_IPS 0 }}
    Right {{ index $VASI_IPS 1 }}

{{- /*
   ---------------------------------------------------------------------------
   * END of template
   ---------------------------------------------------------------------------
*/}}

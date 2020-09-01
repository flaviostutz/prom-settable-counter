package promcollectors

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

//SettableCounter Settable counter for Prometheus with labels
type SettableCounter struct {
	desc          *prometheus.Desc
	labelNames    []string
	labelsCounter map[string]float64
}

//Describe Prometheus will call this to get metrics description
func (c *SettableCounter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

//Collect Prometheus will call this function to collect metric values
func (c *SettableCounter) Collect(ch chan<- prometheus.Metric) {
	for labels, v := range c.labelsCounter {
		lv := strings.Split(labels, "|")[1:]
		ch <- prometheus.MustNewConstMetric(
			c.desc,
			prometheus.CounterValue,
			v,
			lv...,
		)
	}
}

//Set the value for the counter. Should never be less than previous value, only if being reset to zero.
func (c *SettableCounter) Set(value float64, labelValues ...string) error {

	if len(labelValues) != len(c.labelNames) {
		return fmt.Errorf("Incorrect number of arguments in labelValues. Expected: %d", len(c.labelNames))
	}

	//calc key for value in map
	lvv := calcKey(labelValues)

	//if it is an existing value, check for validity
	for k, v := range c.labelsCounter {
		if k == lvv {
			if value != 0 && (value < v || value < 0) {
				return fmt.Errorf("Invalid value for counter for labels %s. New value should be greater then last value and > 0", labelValues)
			}
		}
	}

	c.labelsCounter[lvv] = value
	return nil
}

//NewSettableCounterVec creates a new counter vector with metrics labels
func NewSettableCounterVec(opts prometheus.Opts, labelNames []string) *SettableCounter {
	for _, v := range labelNames {
		if strings.Contains(v, "|") {
			panic("Label name cannot contain char '|'")
		}
	}
	return &SettableCounter{
		desc:          prometheus.NewDesc(opts.Name, opts.Help, labelNames, opts.ConstLabels),
		labelNames:    labelNames,
		labelsCounter: make(map[string]float64),
	}
}

func calcKey(labelValues []string) string {
	lvv := ""
	for _, v := range labelValues {
		lvv = fmt.Sprintf("%s|%s", lvv, v)
	}
	return lvv
}

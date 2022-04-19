import itertools
import json
import os

import requests
from bandit.core import constants, docs_utils, test_properties


def comment_on_pull_request(message):
    """Comments a message within a pull request."""
    token = os.getenv("INPUT_GITHUB_TOKEN")
    if not token:
        print(message)
        return

    if os.getenv("GITHUB_EVENT_NAME") == "pull_request":
        with open(os.getenv("GITHUB_EVENT_PATH")) as json_file:
            event = json.load(json_file)
            headers_dict = {
                "Accept": "application/vnd.github.v3+json",
                "Authorization": f"token {token}",
            }

            request_path = f"https://api.github.com/repos/{event['repository']['full_name']}/issues/{event['number']}/comments"
            requests.post(request_path, headers=headers_dict, json={"body": message})


def generate_metrics_table(manager):
    """Generate a metrics table."""
    bits = []

    confidence_mapping = {
        "UNDEFINED": "Undefined confidence",
        "LOW": "Low confidence",
        "MEDIUM": "Medium confidence",
        "HIGH": "High confidence",
    }
    severity_mapping = {
        "UNDEFINED": "Undefined severity",
        "LOW": "Low severity",
        "MEDIUM": "Medium severity",
        "HIGH": "High severity",
    }

    confusion_matrix = get_confusion_matrix(manager)

    # headline
    bits.append("\n### Severity and confidence confusion matrix\n")

    # table header
    bits.append(f"| | {'|'.join([severity_mapping[c] for c in constants.RANKING])}|")
    bits.append("|:-:|:-:|:-:|:-:|:-:|")

    # create rows
    for constant in constants.RANKING:
        undefined = confusion_matrix[f"{constant}_UNDEFINED"]
        low = confusion_matrix[f"{constant}_LOW"]
        medium = confusion_matrix[f"{constant}_MEDIUM"]
        high = confusion_matrix[f"{constant}_HIGH"]
        bits.append(
            f"|{confidence_mapping[constant]}|{undefined}|{low}|{medium}|{high}|"
        )

    return "\n".join([bit for bit in bits])


def get_issue_description(issue, indent, show_lineno=True, show_code=True, lines=-1):
    """Returns details of issue within a description. This description contains:
    - severity
    - confidence
    - location
    - more info
    - code snippet
    """
    bits = []
    bits.append("<details>")
    bits.append(
        f"<summary><strong>[{issue.test_id}:{issue.test}]</strong> {issue.text}</summary>\n<br>\n"
    )

    bits.append(
        f"|<strong>Severity</strong>| {issue.severity.capitalize()} |\n|:-:|:-:|\n|<strong>Confidence</strong>| {issue.confidence.capitalize()} |"
    )

    bits.append(
        f"|<strong>Location<strong>| {issue.fname}:{issue.lineno if show_lineno else ''}:{''} |"
    )

    bits.append(f"|<strong>More Info<strong>| {docs_utils.get_url(issue.test_id)} |\n")

    if show_code:
        bits.append("<br>\n\n```python")
        bits.extend([indent + line for line in issue.get_code(lines, True).split("\n")])
        bits.append("```\n")

    bits.append("</details>")
    return "\n".join([bit for bit in bits])


def get_confusion_matrix(manager):
    """Returns a confusion matrix with all possible RANKING constants."""
    confusion_matrix = {
        "_".join(combination): 0
        for combination in list(itertools.product(constants.RANKING, constants.RANKING))
    }
    for issue_dict in manager.results:
        issue_severity = issue_dict.as_dict().get("issue_severity", 0)
        issue_confidence = issue_dict.as_dict().get("issue_confidence", 0)
        confusion_matrix[f"{issue_severity}_{issue_confidence}"] += 1

    return confusion_matrix


def get_detailed_results(manager, sev_level, conf_level, lines):
    """Get detailed bandit results for every security issue."""
    bits = []
    issues = manager.get_issue_list(sev_level, conf_level)
    baseline = not isinstance(issues, list)
    candidate_indent = " " * 10

    if not len(issues):
        return "\tNo issues identified."

    for issue in issues:
        # if not a baseline or only one candidate we know the issue of
        if not baseline or len(issues[issue]) == 1:
            bits.append(get_issue_description(issue, "", lines=lines))

        # otherwise show the finding and the candidates
        else:
            bits.append(
                get_issue_description(issue, "", show_lineno=False, show_code=False)
            )

            bits.append("\n-- Candidate Issues --")
            for candidate in issues[issue]:
                bits.append(
                    get_issue_description(candidate, candidate_indent, lines=lines)
                )
                bits.append("\n")
    return "\n".join([bit for bit in bits])


def get_verbose_details(manager):
    """Get verbose details."""
    bits = []
    bits.append(f"Files in scope ({len(manager.files_list)}):")
    tpl = "\t%s (score: {SEVERITY: %i, CONFIDENCE: %i})"
    bits.extend(
        [
            tpl % (item, sum(score["SEVERITY"]), sum(score["CONFIDENCE"]))
            for (item, score) in zip(manager.files_list, manager.scores)
        ]
    )
    bits.append(f"Files excluded ({len(manager.excluded_files)}):")
    bits.extend(["\t%s" % fname for fname in manager.excluded_files])
    return "\n".join([bit for bit in bits])


@test_properties.accepts_baseline
def report(manager, fileobj, sev_level, conf_level, lines: int = -1):
    """Prints discovered issues in the text format

    Note
    ----
    Function template taken from:
    https://bandit.readthedocs.io/en/latest/formatters/index.html#example-formatter

    Parameters
    ----------
    manager : _type_
       A bandit manager object
    sev_level : _type_
        Filtering severity level
    conf_level : _type_
        Filtering confidence level
    lines : int, optional
         Number of lines to report, -1 means all, by default -1
    """
    bits = []
    if manager.results_count(sev_level, conf_level):
        if manager.verbose:
            bits.append(get_verbose_details(manager))

        # add general results
        bits.append("## Bandit results")
        bits.append(
            f"<strong>Total lines of code inspected:</strong> {(manager.metrics.data['_totals']['loc'])}"
        )
        bits.append(
            f"<strong>Total lines skipped (#nosec):</strong> {(manager.metrics.data['_totals']['nosec'])}"
        )
        # add metrics table
        bits.append(generate_metrics_table(manager))

        # add detailed security issue results
        bits.append(
            "<details><summary>📋 Click here to see the all possible security issues</summary>\n"
            "<br>\n"
        )
        bits.append(get_detailed_results(manager, sev_level, conf_level, lines))
        bits.append("</details>")

        result = "\n".join([bit for bit in bits]) + "\n"

        comment_on_pull_request(result)
